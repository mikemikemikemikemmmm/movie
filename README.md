使用golang(gin)+kafka+redis+postgresql  
實現完全事件驅動架構  
1.選取座位後先用redis鎖
2.redis鎖成功後用kafka傳給consumer  
3.consumer消費後，確認可以訂位後再寫上sql  
4.前端訂位後用輪詢查詢sql結果  
![DEMO](https://github.com/mikemikemikemikemmmm/movie/blob/main/movie.gif)
