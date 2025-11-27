import { useStore } from "../store"

export const LoadingComponent = () => {
    const loading = useStore((state) => state.isLoading)

    return (
        <div
            style={{
                position: "fixed",
                top: 0,
                left: 0,
                width: "100vw",
                height: "100vh",
                display: loading ? "flex" : "none",
                justifyContent: "center",
                alignItems: "center",
                zIndex: 99,
                backgroundColor: "rgba(0, 0, 0, 0.7)"
            }}
        >
            <div style={{ padding: 10, backgroundColor: "white", borderRadius: 8 }}>
                訂位處理中
            </div>
        </div>
    )
}