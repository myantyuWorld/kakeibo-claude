import { z } from 'zod'

export const loginSchema = z.object({
  email: z
    .string({ required_error: '必須項目です' })
    .min(1, '必須項目です')
    .email('メールアドレスの形式が正しくありません'),
  password: z.string({ required_error: '必須項目です' }).min(1, '必須項目です'),
})

export type LoginForm = z.infer<typeof loginSchema>
