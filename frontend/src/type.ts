export interface ReserveData {
    seatId: number,
    userId: number
}
export interface Order {
    id: number
    success: boolean
    detail: string
}
export type SeatStatus =  "reserved" | "available" | "broken"
export interface SeatDataType {
    x: number,
    y: number,
    ID: number,
    status:SeatStatus
}
export interface PollingSeatResponse {
    status: "success" | "fail" | "pending"
    detail: string
}
export interface PollingSeatData {
    userId: number
}
export type SeatMatType = (SeatDataType | null)[][]