interface DialogOptions {
  title?: string;
  message: string;
  confirmText?: string;
  cancelText?: string;
  confirmColor?: "red" | "primary" | "green";
}

let open = $state(false);
let options = $state<DialogOptions>({ message: "" });
let resolver: ((value: boolean) => void) | null = null;

function confirm(opts: DialogOptions | string): Promise<boolean> {
  options = typeof opts === "string" ? { message: opts } : opts;
  open = true;

  return new Promise<boolean>((resolve) => {
    resolver = resolve;
  });
}

function handleAction(confirmed: boolean) {
  open = false;
  resolver?.(confirmed);
  resolver = null;
}

export const dialogService = {
  get open() {
    return open;
  },
  set open(value: boolean) {
    open = value;
    if (!value) handleAction(false);
  },
  get options() {
    return options;
  },
  confirm,
  handleAction,
};
