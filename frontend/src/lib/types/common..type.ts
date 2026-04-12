export interface MenuItem {
  label: string;
  onClick?: VoidFunction;
  link?: string;
}

export interface BaseModel {
  id: string;
  createdAt: string;
  createdBy: string;
  updatedAt?: string;
  updatedBy?: string;
}

export interface ListModel<T> {
  data: T[];
  message: string;
}

export interface SuccessModel<T> {
  data: T;
  message: string;
}

export interface ErrorModel {
  message: string;
}
