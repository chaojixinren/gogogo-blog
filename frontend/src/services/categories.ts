import api from './api'
import type { Category, CategoryInput, Paginated, Post } from '@/types'
import type { PostQuery } from './posts'

export const fetchCategories = async (): Promise<Category[]> => {
  const { data } = await api.get<{ data: Category[] }>('/categories')
  return data.data
}

export const createCategory = async (
  payload: CategoryInput,
): Promise<Category> => {
  const { data } = await api.post<{ data: Category }>('/categories', payload)
  return data.data
}

export const updateCategory = async (
  id: number,
  payload: Partial<CategoryInput>,
): Promise<Category> => {
  const { data } = await api.put<{ data: Category }>(
    `/categories/${id}`,
    payload,
  )
  return data.data
}

export const deleteCategory = async (id: number): Promise<void> => {
  await api.delete(`/categories/${id}`)
}

export const fetchPostsByCategory = async (
  categoryIdOrSlug: number | string,
  params: PostQuery = {},
): Promise<Paginated<Post>> => {
  const { data } = await api.get<Paginated<Post>>(
    `/categories/${categoryIdOrSlug}/posts`,
    { params },
  )
  return data
}
