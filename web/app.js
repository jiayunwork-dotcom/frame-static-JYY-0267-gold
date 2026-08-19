"use strict";

const $ = (id) => document.getElementById(id);

async function loadExample() {
  const name = $("example").value;
  const res = await fetch("example/" + name);
  if (!res.ok) {
    setStatus("示例加载失败: " + res.status, true);
    return;
  }
  $("model").value = await res.text();
  setStatus("已加载 " + name);
}

async function solve() {
  const text = $("model").value.trim();
  if (!text) {
    setStatus("请先填入模型 JSON", true);
    return;
  }
  setStatus("计算中…");
  let resp;
  try {
    resp = await fetch("/api/solve", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: text,
    });
  } catch (e) {
    setStatus("请求失败: " + e, true);
    return;
  }
  const data = await resp.json();
  if (!data.ok) {
    setStatus("求解错误: " + (data.error ? data.error.message : "未知"), true);
    return;
  }
  setStatus("求解完成");
  render(data);
}

function setStatus(msg, isError) {
  const el = $("status");
  el.textContent = msg;
  el.className = "status" + (isError ? " error" : "");
}

function render(data) {
  renderReactions(data.reactions || []);
  renderMembers(data.members || []);
  drawFrame(data);
}

function renderReactions(rows) {
  const tb = $("reactions").querySelector("tbody");
  tb.innerHTML = "";
  for (const r of rows) {
    const tr = document.createElement("tr");
    tr.innerHTML = `<td>${r.node}</td><td>${r.dof}</td><td>${r.force.toFixed(4)}</td>`;
    tb.appendChild(tr);
  }
}

function renderMembers(rows) {
  const tb = $("members").querySelector("tbody");
  tb.innerHTML = "";
  for (const m of rows) {
    const tr = document.createElement("tr");
    tr.innerHTML =
      `<td>${m.from}–${m.to}</td>` +
      `<td>${m.Ni.toFixed(3)}</td><td>${m.Vi.toFixed(3)}</td><td>${m.Mi.toFixed(3)}</td>` +
      `<td>${m.Nj.toFixed(3)}</td><td>${m.Vj.toFixed(3)}</td><td>${m.Mj.toFixed(3)}</td>`;
    tb.appendChild(tr);
  }
}

function drawFrame(data) {
  const svg = $("canvas");
  while (svg.firstChild) svg.removeChild(svg.firstChild);
  const nodes = data.nodes || [];
  if (!nodes.length) return;
  const xs = nodes.map((n) => n.x);
  const ys = nodes.map((n) => n.y);
  const minX = Math.min(...xs), maxX = Math.max(...xs);
  const minY = Math.min(...ys), maxY = Math.max(...ys);
  const pad = 50;
  const W = 640, H = 420;
  const sx = (W - 2 * pad) / Math.max(1e-9, maxX - minX);
  const sy = (H - 2 * pad) / Math.max(1e-9, maxY - minY);
  const scale = Math.min(sx, sy);
  const toPx = (x, y) => [
    pad + (x - minX) * scale,
    H - (pad + (y - minY) * scale),
  ];

  // deformation amplification: scale displacement to ~10% of model size
  const dispScale = 0.1 * (maxX - minX) / Math.max(1e-9, maxDisp(nodes));
  const pos = {};
  for (const n of nodes) {
    const [px, py] = toPx(n.x, n.y);
    pos[n.id] = { px, py, dx: n.ux * dispScale, dy: -n.uy * dispScale };
  }

  // original shape (thin)
  const members = data.members || [];
  for (const m of members) {
    const a = pos[m.from], b = pos[m.to];
    if (!a || !b) continue;
    line(svg, a.px, a.py, b.px, b.py, "orig");
  }
  // deformed shape (dashed)
  for (const m of members) {
    const a = pos[m.from], b = pos[m.to];
    if (!a || !b) continue;
    line(svg, a.px + a.dx, a.py + a.dy, b.px + b.dx, b.py + b.dy, "defo");
  }
  // nodes
  for (const n of nodes) {
    const p = pos[n.id];
    circle(svg, p.px + p.dx, p.py + p.dy, n.id);
  }
}

function maxDisp(nodes) {
  let m = 0;
  for (const n of nodes) {
    m = Math.max(m, Math.hypot(n.ux || 0, n.uy || 0));
  }
  return m;
}

function line(svg, x1, y1, x2, y2, cls) {
  const el = document.createElementNS("http://www.w3.org/2000/svg", "line");
  el.setAttribute("x1", x1); el.setAttribute("y1", y1);
  el.setAttribute("x2", x2); el.setAttribute("y2", y2);
  el.setAttribute("class", cls);
  svg.appendChild(el);
}

function circle(svg, x, y, label) {
  const el = document.createElementNS("http://www.w3.org/2000/svg", "circle");
  el.setAttribute("cx", x); el.setAttribute("cy", y); el.setAttribute("r", 4);
  el.setAttribute("class", "node");
  svg.appendChild(el);
  const t = document.createElementNS("http://www.w3.org/2000/svg", "text");
  t.setAttribute("x", x + 6); t.setAttribute("y", y - 6);
  t.setAttribute("class", "nlabel");
  t.textContent = label;
  svg.appendChild(t);
}

$("load").addEventListener("click", loadExample);
$("solve").addEventListener("click", solve);
