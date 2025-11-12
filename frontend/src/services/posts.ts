import api from './api'
import type { Paginated, Post, PostInput } from '@/types'

export interface PostQuery {
  page?: number
  pageSize?: number
  status?: string
  tag?: string
  category?: string
  author?: string
  search?: string
  includeContent?: boolean
}

export const fetchPosts = async (
  params: PostQuery = {},
): Promise<Paginated<Post>> => {
  const { data } = await api.get<Paginated<Post>>('/posts', { params })
  return data
}

export const fetchPostBySlug = async (slug: string): Promise<Post> => {
  const { data } = await api.get<{ data: Post }>(`/posts/slug/${slug}`)
  return data.data
}

export const fetchPostById = async (id: number): Promise<Post> => {
  const { data } = await api.get<{ data: Post }>(`/posts/${id}`)
  return data.data
}

export const fetchMyPosts = async (
  params: PostQuery = {},
): Promise<Paginated<Post>> => {
  const { data } = await api.get<Paginated<Post>>('/me/posts', { params })
  return data
}

export const createPost = async (payload: PostInput): Promise<Post> => {
  const { data } = await api.post<{ data: Post }>('/posts', payload)
  return data.data
}

export const updatePost = async (
  id: number,
  payload: Partial<PostInput>,
): Promise<Post> => {
  const { data } = await api.put<{ data: Post }>(`/posts/${id}`, payload)
  return data.data
}

export const deletePost = async (id: number): Promise<void> => {
  await api.delete(`/posts/${id}`)
}
