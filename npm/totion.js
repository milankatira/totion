#!/usr/bin/env node
// totion npm shim — downloads the Go binary for this platform from GitHub
// Releases on first run, then executes it. Dependency-free.

"use strict";

const fs = require("fs");
const path = require("path");
const https = require("https");
const { spawnSync } = require("child_process");

const pkg = require("./package.json");
const REPO = "milankatira/totion";

function targetPlatform() {
  const os = { darwin: "darwin", linux: "linux", win32: "windows" }[process.platform];
  const arch = { x64: "amd64", arm64: "arm64" }[process.arch];
  if (!os || !arch) {
    console.error(`totion: unsupported platform: ${process.platform}/${process.arch}`);
    console.error(`Download a binary manually from https://github.com/${REPO}/releases/latest`);
    process.exit(1);
  }
  return { os, arch, ext: os === "windows" ? ".exe" : "" };
}

function download(url, dest, redirectsLeft = 5) {
  return new Promise((resolve, reject) => {
    https
      .get(url, { headers: { "User-Agent": `totion-npm/${pkg.version}` } }, (res) => {
        if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
          res.resume();
          if (redirectsLeft === 0) return reject(new Error("too many redirects"));
          return resolve(download(res.headers.location, dest, redirectsLeft - 1));
        }
        if (res.statusCode !== 200) {
          res.resume();
          return reject(new Error(`download failed: HTTP ${res.statusCode} for ${url}`));
        }
        const file = fs.createWriteStream(dest, { mode: 0o755 });
        res.pipe(file);
        file.on("finish", () => file.close(resolve));
        file.on("error", reject);
      })
      .on("error", reject);
  });
}

async function ensureBinary() {
  const { os, arch, ext } = targetPlatform();
  const binDir = path.join(__dirname, "bin");
  const binPath = path.join(binDir, `totion${ext}`);

  if (fs.existsSync(binPath)) return binPath;

  // Pin the binary to this package's version so npm version == binary version.
  const url = `https://github.com/${REPO}/releases/download/v${pkg.version}/totion_${os}_${arch}${ext}`;

  fs.mkdirSync(binDir, { recursive: true });
  process.stderr.write(`totion: downloading v${pkg.version} (${os}/${arch})...\n`);
  try {
    await download(url, binPath);
  } catch (err) {
    try {
      fs.unlinkSync(binPath);
    } catch {
      // best-effort cleanup of a partial download
    }
    console.error(`totion: ${err.message}`);
    console.error(`Try installing another way: https://github.com/${REPO}#install`);
    process.exit(1);
  }
  return binPath;
}

async function main() {
  const binPath = await ensureBinary();
  const result = spawnSync(binPath, process.argv.slice(2), { stdio: "inherit" });
  if (result.error) {
    console.error(`totion: failed to run binary: ${result.error.message}`);
    process.exit(1);
  }
  process.exit(result.status ?? 0);
}

main();
