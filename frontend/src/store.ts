import { create } from 'zustand'

interface StoreState {
    isLoading: boolean;
    setLoading: (p: boolean) => void;
}
export const useStore = create<StoreState>((set) => ({
    isLoading: false,
    setLoading: (payload: boolean) =>
        set(() => ({ isLoading: payload })),
}));