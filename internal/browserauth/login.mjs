import fs from 'node:fs';
import { chromium } from 'playwright';

const outputPath = process.env.FOODCLI_OUTPUT_PATH || process.env.FOODORACLI_OUTPUT_PATH;
if (!outputPath) {
  process.stderr.write('FOODCLI_OUTPUT_PATH missing\n');
  process.exit(2);
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function readStdinJSON() {
  const chunks = [];
  for await (const chunk of process.stdin) chunks.push(chunk);
  const raw = Buffer.concat(chunks).toString('utf8').trim();
  if (!raw) throw new Error('stdin empty');
  return JSON.parse(raw);
}

function oauthURL(baseURL) {
  const u = new URL(baseURL);
  if (!u.pathname.endsWith('/')) u.pathname += '/';
  u.pathname += 'oauth2/token';
  return u.toString();
}

function isHTML(status, headers, body) {
  if (status !== 403 && status !== 429 && status !== 503) return false;
  const ct = (headers['content-type'] || '').toLowerCase();
  if (ct.includes('text/html')) return true;
  const b = (body || '').trimStart();
  return b.startsWith('<!DOCTYPE html') || b.startsWith('<html');
}

function tryParseJSON(body) {
  try {
    return JSON.parse(body);
  } catch {
    return null;
  }
}

function isPerimeterXBlocked(status, headers, body) {
  if (status !== 403) return false;
  const ct = (headers['content-type'] || '').toLowerCase();
  if (!ct.includes('application/json')) return false;
  const obj = tryParseJSON(body);
  return !!(obj && (obj.appId || obj.app_id) && (obj.blockScript || obj.altBlockScript));
}

function perimeterXBlockURL(baseURL, body) {
  const obj = tryParseJSON(body);
  if (!obj) return '';
  const baseOrigin = new URL(baseURL).origin;
  const rel = obj.blockScript;
  if (typeof rel === 'string' && rel.startsWith('/')) {
    return new URL(rel, baseOrigin).toString();
  }
  const alt = obj.altBlockScript;
  if (typeof alt === 'string' && alt.startsWith('http')) {
    return alt;
  }
  return '';
}

const input = await readStdinJSON();
const timeoutMillis = Math.max(10_000, Number(input.timeout_millis || 0));
const deadline = Date.now() + timeoutMillis;

const browser = await chromium.launch({ headless: false });
const context = await browser.newContext();
const page = await context.newPage();

try {
  // Intentional: let the user complete any Cloudflare/PerimeterX checks in a real browser.
  // Use origin (not /api/v5/) so challenge pages can render.
  await page.goto(new URL(input.base_url).origin, { waitUntil: 'domcontentloaded' }).catch(() => {});

  const url = oauthURL(input.base_url);

  let lastLog = 0;
  // Loop until the oauth call stops returning Cloudflare HTML (user solved the check), or timeout.
  while (Date.now() < deadline) {
    const res = await context.request.post(url, {
      form: {
        username: input.email,
        password: input.password,
        grant_type: 'password',
        client_secret: input.client_secret,
        scope: 'API_CUSTOMER',
        client_id: input.client_id || 'android',
      },
      headers: {
        Accept: 'application/json',
        'X-Device': input.device_id,
        'X-OTP-Method': input.otp_method || 'sms',
        ...(input.otp_code ? { 'X-OTP': input.otp_code } : {}),
        ...(input.mfa_token ? { 'X-Mfa-Token': input.mfa_token } : {}),
      },
    });

    const status = res.status();
    const headers = res.headers();
    const body = await res.text();

    if (isHTML(status, headers, body) || isPerimeterXBlocked(status, headers, body)) {
      if (Date.now() - lastLog > 5000) {
        lastLog = Date.now();
        process.stderr.write('waiting for browser clearance (solve the challenge in the opened window)...\n');
      }
      const pxURL = perimeterXBlockURL(input.base_url, body);
      if (pxURL) {
        await page.goto(pxURL, { waitUntil: 'domcontentloaded' }).catch(() => {});
      }
      await sleep(1500);
      continue;
    }

    const cookies = await context.cookies(new URL(input.base_url).origin);
    const cookieHeader = cookies.map((c) => `${c.name}=${c.value}`).join('; ');
    const userAgent = await page.evaluate(() => navigator.userAgent).catch(() => '');

    fs.writeFileSync(
      outputPath,
      JSON.stringify({
        status,
        body,
        headers,
        cookie_header: cookieHeader,
        user_agent: userAgent,
      }),
      'utf8',
    );
    await browser.close();
    process.exit(0);
  }

  process.stderr.write('timeout waiting for browser clearance\n');
  await browser.close();
  process.exit(3);
} catch (e) {
  try {
    await browser.close();
  } catch {}
  process.stderr.write(String(e?.stack || e) + '\n');
  process.exit(1);
}
