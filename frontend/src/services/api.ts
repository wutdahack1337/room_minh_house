import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8001';

const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 5000,
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface HealthResponse {
  status: string;
}

export const getRooms = async () => {
  const response = await apiClient.get('/api/rooms');
  return response.data?.rooms ?? response.data ?? [];
};

export const getUsers = async () => {
  const response = await apiClient.get('/api/users');
  return response.data?.users ?? response.data ?? [];
};

export const getContracts = async () => {
  const response = await apiClient.get('/api/contracts');
  return response.data?.contracts ?? response.data ?? [];
};

export const getHealth = async (): Promise<HealthResponse> => {
  const response = await apiClient.get<HealthResponse>('/api/health');
  return response.data;
};