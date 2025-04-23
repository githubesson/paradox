import { authFetch } from './auth';


const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export interface Log {
  uuid: string;
  build_id: string;
  timestamp: string;
  path: string;
  computer_name: string;
  user_name: string;
  country_name: string;
  city_name: string;
}

export interface Build {
  build_id: string;
  filename: string;
  timestamp: string;
}

export interface LogsResponse {
  logs: Log[];
  count: number;
}

export interface BuildsResponse {
  builds: Build[];
  count: number;
}

export interface BuildResponse {
  message: string;
  build_id: string;
  filename: string;
}


export async function fetchLogs() {
  try {
    const response = await authFetch(`${API_BASE_URL}/logs`);
    if (!response.ok) {
      throw new Error(`Failed to fetch logs: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error("Error fetching logs:", error);
    throw error;
  }
}


export async function fetchBuilds() {
  try {
    const response = await authFetch(`${API_BASE_URL}/builds`);
    if (!response.ok) {
      throw new Error(`Failed to fetch builds: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error("Error fetching builds:", error);
    throw error;
  }
}


export async function triggerBuild() {
  try {
    const response = await authFetch(`${API_BASE_URL}/build`);
    if (!response.ok) {
      throw new Error(`Failed to trigger build: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error("Error triggering build:", error);
    throw error;
  }
}


export async function downloadBuild(buildId: string) {
  try {
    const response = await authFetch(`${API_BASE_URL}/download/build/${buildId}`);
    if (!response.ok) {
      throw new Error(`Failed to download build: ${response.status}`);
    }

    
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `build-${buildId}`;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);

    return true;
  } catch (error) {
    console.error("Error downloading build:", error);
    throw error;
  }
}


export async function downloadLogs(uuid: string) {
  try {
    const response = await authFetch(`${API_BASE_URL}/download/logs/${uuid}`);
    if (!response.ok) {
      throw new Error(`Failed to download logs: ${response.status}`);
    }

    
    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `logs-${uuid}.zip`;
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);

    return true;
  } catch (error) {
    console.error("Error downloading logs:", error);
    throw error;
  }
}
