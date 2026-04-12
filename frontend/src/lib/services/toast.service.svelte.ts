type ToastColor = "green" | "red" | "yellow" | "blue";

export interface ToastItem {
  id: number;
  message: string;
  color: ToastColor;
  timeoutId?: ReturnType<typeof setTimeout>;
}

const toasts = $state<ToastItem[]>([]);
let nextId = 0;

function add(message: string, color: ToastColor, duration = 5000) {
  const id = ++nextId;
  const timeoutId = setTimeout(() => dismiss(id), duration);
  toasts.push({ id, message, color, timeoutId });
}

function dismiss(id: number) {
  const idx = toasts.findIndex((t) => t.id === id);
  if (idx === -1) return;
  clearTimeout(toasts[idx].timeoutId);
  toasts.splice(idx, 1);
}

export const toastService = {
  get items() {
    return toasts;
  },
  success: (msg: string) => add(msg, "green"),
  error: (msg: string) => add(msg, "red"),
  warning: (msg: string) => add(msg, "yellow"),
  info: (msg: string) => add(msg, "blue"),
  dismiss,
};
