import api from './api'
import type { Paginated, Post, Tag, TagInput } from '@/types'
import type { PostQuery } from './posts'

export const fetchTags = async (): Promise<Tag[]> => {
  const { data } = await api.get<{ data: Tag[] }>('/tags')
  return data.data
}

export const createTag = async (payload: TagInput): Promise<Tag> => {
  const { data } = await api.post<{ data: Tag }>('/tags', payload)
  return data.data
}

export const updateTag = async (
  id: number,
  payload: Partial<TagInput>,
): Promise<Tag> => {
  const { data } = await api.put<{ data: Tag }>(`/tags/${id}`, payload)
  return data.data
}

export const deleteTag = async (id: number): Promise<void> => {
  await api.delete(`/tags/${id}`)
}

export const fetchPostsByTag = async (
  slug: string,
  params: PostQuery = {},
): Promise<Paginated<Post>> => {
  const { data } = await api.get<Paginated<Post>>(`/tags/${slug}/posts`, {
    params,
  })
  return data
}
