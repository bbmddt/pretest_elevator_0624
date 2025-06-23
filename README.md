# 電梯管理系統 Bollertech Interview 0624

## 系統設計

### 本系統採用混合型調度策略：

1. **SCAN（電梯掃描算法）特徵**
   - 電梯保持單一方向運行直到該方向沒有請求
   - 分別維護向上和向下的請求隊列

2. **SSTF（最短尋找時間優先）特徵**
   - 優先選擇最近的待服務樓層
   - 通過距離計算最小化移動時間

3. **優先級處理**
   - 電梯內乘客優先於等待乘客
   - 方向相同的請求優先處理

### 系統運作流程

**初始化階段**
   
- 初始化電梯，創建 Building 物件，建立 Dispatcher

**運作循環**
- 乘客生成
- 電梯調度分配
- 電梯移動與乘客處理
- 時間更新

**結束條件**
- 所有乘客都已被生成
- 所有電梯都處於閒置狀態
- 所有乘客都已到達目的地

### 並發控制

- 使用互斥鎖確保資源安全訪問
- Building 和 Elevator 各自維護獨立的鎖
- Dispatcher 協調各個元件的運作

```
Dispatcher
    ├── Building Lock (保護建築物共享資源)
    │   ├── WaitQueue
    │   ├── PeopleDone
    │   └── TimeElapsed
    │
    └── Elevator Locks (保護各電梯獨立狀態)
        ├── Elevator 1 (CurrentFloor, Direction, Passengers)
        └── Elevator 2 (CurrentFloor, Direction, Passengers)
```

## 主要函式

#### Building
```
NewBuilding(): 初始化建築物，設置樓層數、電梯數量等
GeneratePerson(): 產生新乘客並加入等待隊列
```
#### Elevator
```
NewElevator(): 初始化電梯，設定初始樓層等
Step(): 執行電梯一個時間單位的動作
updateDirection(): 更新電梯運行方向
```
#### Dispatcher
```
Run(): 控制整體模擬流程
getWaitingFloors(): 獲取有待乘客的樓層
isSimulationDone(): 檢查模擬是否完成
```

## 物件結構

#### Person
```go
type Person struct {
    ID          int    // 乘客唯一識別碼
    Origin      int    // 起始樓層
    Destination int    // 目標樓層
}
```

#### Building
```go
type Building struct {
    Floors      int    // 大樓總樓層數
    TotalPeople int    // 總乘客數
    WaitQueue   map[int]struct {  // 每層樓的等待隊列
        UpQueue   []*Person       // 向上等待隊列
        DownQueue []*Person       // 向下等待隊列
    }
    Elevators   []*Elevator      // 電梯列表
    PeopleDone  int             // 已完成運送的乘客數
    TimeElapsed int             // 已耗用時間
    Mutex       sync.Mutex      // 同步鎖
}
```

#### Elevator
```go
type Elevator struct {
    ID           int       // 電梯編號
    CurrentFloor int       // 當前樓層
    Direction    string    // 移動方向："up", "down", "idle"
    Passengers   []*Person // 電梯內的乘客
    Mutex        sync.Mutex // 同步鎖
}
```

#### Dispatcher
```go
type Dispatcher struct {
    Building *Building    // 對應的建築物
}
```
