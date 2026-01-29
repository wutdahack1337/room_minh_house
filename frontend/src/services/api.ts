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

export const getHealth = async (): Promise<HealthResponse> => {
    const response = await apiClient.get<HealthResponse>('/api/health');
    return response.data;
};