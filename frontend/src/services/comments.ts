import api from './api'
import type { Comment } from '@/types'

export const fetchComments = async (
  postIdOrSlug: number | string,
): Promise<Comment[]> => {
  const { data } = await api.get<{ data: Comment[] }>(
    `/posts/${postIdOrSlug}/comments`,
  )
  return data.data
}

export const createComment = async (
  postIdOrSlug: number | string,
  payload: { authorName?: string; body: string },
): Promise<Comment> => {
  const { data } = await api.post<{ data: Comment }>(
    `/posts/${postIdOrSlug}/comments`,
    payload,
  )
  return data.data
}
