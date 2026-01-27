const LOG_LEVEL = process.env.LOG_LEVEL || "error";
var LogLevel = /* @__PURE__ */ ((LogLevel2) => {
  LogLevel2[LogLevel2["off"] = 0] = "off";
  LogLevel2[LogLevel2["info"] = 1] = "info";
  LogLevel2[LogLevel2["warn"] = 2] = "warn";
  LogLevel2[LogLevel2["error"] = 3] = "error";
  return LogLevel2;
})(LogLevel || {});
const logger = {
  info: (msg, meta) => 1 >= LogLevel[LOG_LEVEL] ? console.log(JSON.stringify({ level: "info", msg, ...meta, timestamp: Date.now() })) : null,
  warn: (msg, meta) => 2 >= LogLevel[LOG_LEVEL] ? console.warn(JSON.stringify({ level: "warn", msg, ...meta, timestamp: Date.now() })) : null,
  error: (msg, error) => console.error(JSON.stringify({ level: "error", msg, error: error?.message, timestamp: Date.now() }))
};
export {
  logger as l
};
