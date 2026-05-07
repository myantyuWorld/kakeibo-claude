import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { type LoginForm, loginSchema } from './schemas'

export const useLoginForm = () => {
  const { defineField, handleSubmit, errors } = useForm<LoginForm>({
    validationSchema: toTypedSchema(loginSchema),
  })

  const onSubmit = handleSubmit(() => {
    // TODO: 認証方法確定後に API コールを実装する
  })

  return { defineField, errors, onSubmit }
}
