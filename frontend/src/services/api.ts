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

export const listRooms = async () => {
  const response = await apiClient.get('/api/rooms');
  return response.data?.rooms ?? response.data ?? [];
};
export const createRoom = async (data: { room_number: number; price: number; is_available: boolean }) => {
  const r = await apiClient.post('/api/rooms', data);
  return r.data;
};

export const listUsers = async () => {
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

export const listContracts = async () => {
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

export const listInvoices = async (contractId: string) => {
  const r = await apiClient.get(`/api/contracts/${contractId}/invoices`);
  return r.data?.invoices ?? r.data ?? [];
};
export const createInvoice = async (data: {
  contract_id: string;
  period: string;
  electricity_unit_price: number;
  water_unit_price: number;
}) => {
  const r = await apiClient.post(`/api/contracts/${data.contract_id}/invoices`, data);
  return r.data;
};

export const listInvoiceItems = async (contractId: string, period: string) => {
  const r = await apiClient.get(`/api/contracts/${contractId}/invoices/${period}/items`);
  return r.data?.invoice_items ?? r.data ?? [];
};

export const updateInvoiceItem = async (item: {
  contract_id: string;
  period: string;
  item_type: string;
  unit_price: number;
  quantity: number;
  amount: number;
}) => {
  const r = await apiClient.put(
    `/api/contracts/${item.contract_id}/invoices/${item.period}/items/${item.item_type}`,
    item
  );
  return r.data;
};

export const getHealth = async (): Promise<HealthResponse> => {
  const response = await apiClient.get<HealthResponse>('/api/health');
  return response.data;
};