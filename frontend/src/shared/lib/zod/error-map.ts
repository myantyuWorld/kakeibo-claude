import { type ZodErrorMap, ZodIssueCode, ZodParsedType } from 'zod'

export const zodErrorMapJa: ZodErrorMap = (issue, ctx) => {
  switch (issue.code) {
    case ZodIssueCode.invalid_type:
      if (issue.received === ZodParsedType.undefined) {
        return { message: '必須項目です' }
      }
      return { message: `${issue.expected}型で入力してください` }

    case ZodIssueCode.too_small:
      if (issue.type === 'string') {
        if (issue.minimum === 1) return { message: '必須項目です' }
        return { message: `${issue.minimum}文字以上で入力してください` }
      }
      if (issue.type === 'number') {
        return { message: `${issue.minimum}以上の値を入力してください` }
      }
      return { message: ctx.defaultError }

    case ZodIssueCode.too_big:
      if (issue.type === 'string') {
        return { message: `${issue.maximum}文字以内で入力してください` }
      }
      if (issue.type === 'number') {
        return { message: `${issue.maximum}以下の値を入力してください` }
      }
      return { message: ctx.defaultError }

    case ZodIssueCode.invalid_string:
      if (issue.validation === 'email') {
        return { message: 'メールアドレスの形式が正しくありません' }
      }
      if (issue.validation === 'url') {
        return { message: 'URLの形式が正しくありません' }
      }
      return { message: ctx.defaultError }

    case ZodIssueCode.invalid_enum_value:
      return { message: '選択肢から選んでください' }

    default:
      return { message: ctx.defaultError }
  }
}
