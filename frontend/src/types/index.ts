export interface User {
  id: number
  username: string
  email?: string
  displayName: string
  bio?: string
  avatarUrl?: string
  createdAt: string
}

export interface Category {
  id: number
  name: string
  slug: string
  description?: string
  createdAt: string
}

export interface Tag {
  id: number
  name: string
  slug: string
  createdAt: string
}

export interface Comment {
  id: number
  authorName: string
  body: string
  approved: boolean
  createdAt: string
  user?: User
}

export interface Post {
  id: number
  title: string
  summary?: string
  content: string
  slug: string
  status: 'draft' | 'published' | 'archived'
  coverImage?: string
  publishedAt?: string | null
  author: User
  category?: Category | null
  tags: Tag[]
  comments?: Comment[]
  createdAt: string
  updatedAt: string
}

export interface Paginated<T> {
  data: T[]
  page: number
  pageSize: number
  total: number
}

export interface PostInput {
  title: string
  summary?: string
  content: string
  status: string
  slug?: string
  categoryId?: number | null
  categorySlug?: string
  tags?: string[]
  coverImage?: string
  publishedAt?: string | null
}

export interface CategoryInput {
  name: string
  slug?: string
  description?: string
}

export interface TagInput {
  name: string
  slug?: string
}

