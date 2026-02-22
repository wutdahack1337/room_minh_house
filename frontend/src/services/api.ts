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
export const createRoom = async (data: { room_number: number; price: number; is_available: boolean }) => {
  const r = await apiClient.post('/api/rooms', data);
  return r.data;
};

export const getUsers = async () => {
  const response = await apiClient.get('/api/users');
  return response.data?.users ?? response.data ?? [];
};
export const createUser = async (data: {
  pid: string;
  phone_number: string;
  role: string;
  id_card: { name: string; sex: string; date_of_birth: string; place_of_birth: string };
  account: { username: string; password: string };
}) => {
  const r = await apiClient.post('/api/users', data);
  return r.data;
};

export const getContracts = async () => {
  const response = await apiClient.get('/api/contracts');
  return response.data?.contracts ?? response.data ?? [];
};
export const createContract = async (data: {
  user_pid: string;
  room_number: number;
  start_date: string;
  duration_months: number;
}) => {
  const r = await apiClient.post('/api/contracts', data);
  return r.data;
};

export const getHealth = async (): Promise<HealthResponse> => {
  const response = await apiClient.get<HealthResponse>('/api/health');
  return response.data;
};