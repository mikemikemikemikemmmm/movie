import React from "react";
import type { SeatDataType } from "../type";
import { SEAT_WIDTH } from "../config";
const Seat = (props: { seatData: SeatDataType, picked: boolean }) => {
    const { status, ID, x, y } = props.seatData
    const bg = () => {
        if (props.picked) {
            return "green"
        }
        if (status !== "available") {
            return "red"
        }
        return "blue"
    }
    return (
        <span className="hover-pointer" style={{
            position: "absolute",
            width: SEAT_WIDTH,
            height: SEAT_WIDTH,
            backgroundColor: bg(),
            top: 2 * y * SEAT_WIDTH - SEAT_WIDTH,
            left: 2 * x * SEAT_WIDTH - SEAT_WIDTH,
        }} data-seatid={ID} data-seatstatus={status}>
        </span>
    );
}
export const MemoSeat = React.memo(Seat, (prev, next) => prev.seatData.status === next.seatData.status && prev.picked === next.picked);
