import { useEffect, useRef, useState } from 'react'
import { getApi, postApi } from './api'
import type { SeatDataType, SeatStatus } from './type'
import { USER_ID } from './config';
import { useStore } from './store';

const maxRetry = 60
export function useSeats() {
    const setLoading = useStore((state) => state.setLoading)
    const [seats, setSeats] = useState<SeatDataType[]>([])
    const [pickedSeats, setPickedSeats] = useState<SeatDataType[]>([])
    const maxX = useRef(0)
    const maxY = useRef(0)
    const retryCount = useRef(0)
    const pollingReserveStatus = async () => {
        if (retryCount.current >= maxRetry) {
            alert("訂位超時，請稍後再試")
            return
        }

        setTimeout(async () => {
            try {
                const { httpCode } = await postApi<string>("check_reserve", {
                    user_id: USER_ID,
                    seat_ids: pickedSeats.map(s => s.ID)
                })
                if (httpCode === 200) {
                    alert("訂位成功")
                    await refresh()
                    setLoading(false)
                    retryCount.current = 0
                } else if (httpCode === 202) {
                    console.log("訂位處理中")
                    retryCount.current++
                    pollingReserveStatus()
                } else {
                    alert("訂位失敗")
                    setLoading(false)
                    retryCount.current = 0
                }
            } catch (e) {
                console.error(e)
                setLoading(false)
            }
        }, 1000)
    }
    const setMax = (seatsList: SeatDataType[]) => {
        let _maxY = -Infinity;
        let _maxX = -Infinity;
        for (const p of seatsList) {
            if (p.y > _maxY) _maxY = p.y;
            if (p.x > _maxX) _maxX = p.x;
        }
        maxX.current = _maxX
        maxY.current = _maxY
    }
    const cleanPicked = () => {
        setPickedSeats([])
    }
    const onclick =
        (e: React.MouseEvent<HTMLDivElement>) => {
            const seatDom = e.target as HTMLDivElement
            const seatId = seatDom.getAttribute("data-seatid");
            if (!seatId) return;

            // const seatStatus = seatDom.getAttribute("data-seatstatus") as null | SeatStatus;
            // if (!seatStatus || seatStatus === "broken" || seatStatus === "reserved") return;
            const targetSeat = seats.find(s => s.ID === Number(seatId))
            if (!targetSeat) {
                return
            }
            setPickedSeats(prev => {
                const copyPicked = [...prev];
                const index = copyPicked.findIndex(p => p.ID === targetSeat.ID)
                if (index !== -1) {
                    copyPicked.splice(index, 1);
                } else {
                    // 不在 picked，加入它
                    copyPicked.push(targetSeat);
                }
                return copyPicked;
            });
        }
    const onsubmit = async () => {
        if (pickedSeats.length === 0) {
            return
        }
        setLoading(true)
        const { error } = await postApi("reserve", {
            user_id: USER_ID,
            seat_ids: pickedSeats.map(s => s.ID)
        })
        cleanPicked()
        if (error) {
            alert(error)
            setLoading(false)
        } else {
            pollingReserveStatus()
        }
    }
    const refresh = async () => {
        const result = await getApi<number[]>("refresh_seats")
        if (result.error) {
            alert("更新錯誤")
            return //TODO
        }
        const reservedSeatIds = result.data
        const copySeats = seats.map(s => {
            if (reservedSeatIds.includes(s.ID)) {
                return { ...s, status: "reserved" as SeatStatus }
            } else {
                return s
            }
        })
        setSeats(copySeats)
    }
    const getAllSeat = async () => {
        const seatsList = await getApi<SeatDataType[]>("seats")
        if (seatsList?.error) {
            alert("錯誤")
            return
        }
        if (seatsList.data.length === 0) {
            return
        }
        setMax(seatsList.data)
        setSeats(seatsList.data)
    }
    useEffect(() => {
        getAllSeat()
    }, [])
    return {
        onclick, onsubmit, seats, maxX, maxY, refresh, pickedSeats
    }
}

