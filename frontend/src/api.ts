import { API_BASE_URL } from "./config";

export interface Resp<T> {
  data: T;
  error?: string; // error 可選
  httpCode: number
}

export const getApi = async <T>(url: string): Promise<Resp<T>> => {
  try {
    const response = await fetch(`${API_BASE_URL}/${url}`);

    if (!response.ok) {
      console.error(`HTTP error! status: ${response.status}`);
      return { data: {} as T, error: `HTTP ${response.status}`, httpCode: response.status };
    }

    // 嘗試解析 JSON
    const json: { data: T } = await response.json();
    const data = json.data
    // 保證返回符合 Resp<T> 結構
    return { data, httpCode: response.status };
  } catch (err: unknown) {
    console.error('Fetch error:', err);
    return { data: {} as T, error: err instanceof Error ? err.message : String(err), httpCode: 400 };
  }
};
export const postApi = async <T>(url: string, dataObj: object): Promise<Resp<T>> => {
  try {
    const response = await fetch(`${API_BASE_URL}/${url}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(dataObj),
    });


    // 嘗試解析 JSON
    const json: { data: T, error: string } = await response.json();

    if (!response.ok) {
      console.error(`HTTP error! status: ${response.status}`);
      return { data: {} as T, error: json.error , httpCode: response.status };
    }

    return { data:json.data, httpCode: response.status };
  } catch (err: unknown) {
    console.error('Fetch error:', err);
    return { data: {} as T, error: err instanceof Error ? err.message : String(err), httpCode: 400 };
  }
};