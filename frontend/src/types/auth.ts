export type User = {
  id: number
  email: string
  created_at: string
  updated_at: string
}

export type AuthCredentials = {
  email: string
  password: string
}

export type RegisterResponse = {
  message: string
  user: User
}

export type LoginResponse = {
  message: string
  token: string
}
