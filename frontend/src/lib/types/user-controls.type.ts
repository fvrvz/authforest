import type { MenuItem } from "./common..type";

export type UserControlsProps = {
  fullName: string;
  email: string;
  menuItems: MenuItem[];
  signoutFn: VoidFunction;
};
