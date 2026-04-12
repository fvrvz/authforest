import dayjs from "dayjs";
import z from "zod/v4";

export const userSchema = z.object({
  firstName: z.string().trim(),
  lastName: z.string().optional(),
  email: z.email(),
  username: z.string().trim(),
  password: z.string().trim(),
  DOB: z.coerce.date().transform((date) => dayjs(date).format("YYYY-MM-DD")),
});
