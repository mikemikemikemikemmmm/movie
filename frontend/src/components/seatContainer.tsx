import { MemoSeat } from './seat';
import { SEAT_WIDTH } from '../config';
import { useSeats } from '../hook';

export function SeatsContainer() {
    const { maxX, maxY, onclick, onsubmit, seats, pickedSeats, refresh } = useSeats()
    return (
        <div style={{ display: "flex", justifyContent: "start" }}>
            <div style={{ width: "50%", display: "inline-block", padding: 20 }}>
                <div style={{ height: 20 }}> </div>
                <div>
                    <div style={{ minWidth: maxX.current * SEAT_WIDTH * 2 + SEAT_WIDTH, display: 'inline-flex', justifyContent: "center", alignItems: "center", border: "1px solid black", padding: 5 }}>
                        電影螢幕
                    </div>
                </div>
                <div onClick={onclick} style={{
                    display: "inline-block",
                    position: "relative",
                    minWidth: maxX.current * SEAT_WIDTH * 2 + SEAT_WIDTH,
                    minHeight: maxY.current * SEAT_WIDTH * 2 + SEAT_WIDTH
                }}>
                    {
                        seats.map(s => <MemoSeat
                            key={s.ID}
                            seatData={s}
                            picked={pickedSeats.findIndex(p => p.ID === s.ID) !== -1}
                        />)
                    }
                </div>
            </div>
            <div style={{ width: "50%", display: "inline-block", padding: 20 }}>
                <div >
                    <div style={{ textAlign: "center" }}>
                        已選擇座位
                    </div>
                    {
                        pickedSeats.map(s => <div key={s.ID}>
                            {`${s.y}排${s.x}座`}
                        </div>)
                    }
                </div>
                <div>
                    <button onClick={onsubmit}>
                        確認訂位
                    </button>
                    <button onClick={refresh}>
                        更新座位狀況
                    </button>
                </div>
            </div>
        </div>
    )
}

