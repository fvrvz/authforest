import dayjs from "dayjs";

export function parseBackendDate(date: string) {
  // normalize microseconds → milliseconds
  const normalized = date.replace(/(\.\d{3})\d+/, "$1");
  return dayjs(normalized);
}

export function formatBackendDate(date: string, format = "DD-MM-YYYY") {
  const d = parseBackendDate(date);
  return d.isValid() ? d.format(format) : "";
}
