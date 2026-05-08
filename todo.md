# WoolfCLI 分階段 TODO

本 TODO 依據 `docs/woolf-spec.md` 與 `docs/testing.md` 整理，目標是先完成 Phase 1 Must Have，再逐步補齊 Should Have 與 Phase 2 功能。每一階段都應維持現有 Go 模組邊界：CLI/TUI 只負責互動，Orchestrator 負責流程，Agents/OpenRouter/Ingestion/Session/Exporter/Cost/Config 各自封裝。

## Phase 0：專案基礎與現有骨架校準

- [ ] 對照 spec 的模組總覽，確認 `internal/*` 與 `pkg/pdfparse` 目前檔案職責是否一致。
- [ ] 補齊 CLI command surface：`delete`、`fork`、`agents`、`version`，並確認既有 `init/start/resume/list/show/export/config/models` 行為符合 spec。
- [ ] 統一錯誤碼與錯誤分類，先建立可被 CLI/TUI 顯示的錯誤型別。
- [ ] 建立或補齊測試資料：`testdata/sample.md`、`sample.txt`、`sample.pdf`、session fixture。
- [ ] 確認 `scripts/test.ps1`、Makefile、GitHub Actions 與 `docs/testing.md` 一致。
- [ ] 為 root CLI smoke test 補上所有 Phase 1 指令檢查，避免指令表面回歸。

## Phase 1A：設定、Agent 與 Preset 系統

- [ ] 完成 TOML 設定載入與預設值合併，支援 `--config` 覆蓋預設路徑。（M-15）
- [ ] 支援 OpenRouter API key 從設定檔與環境變數載入，所有輸出必須遮蔽敏感值。（M-14、M-15）
- [ ] 完成 6 個內建 Agent 角色：strict-editor、casual-reader、structure-analyst、marketing-eye、advocate、challenger。（M-03）
- [ ] 完成 YAML 自訂 Agent 載入、驗證與錯誤回報。（M-03）
- [ ] 完成至少 3 組 preset，例如 editorial、brainstorm、reader-review。（M-04）
- [ ] 完成 prompt 模板組裝，包含角色 system prompt、稿件、前序回應、stance 指示與使用者介入內容。（M-01、M-02、M-16）
- [ ] 補齊 `woolf agents list/show/add/edit/delete` 與 `woolf agents preset list/show` 的基本行為。（M-03、M-04）
- [ ] 測試：Agent YAML 驗證、preset 展開、prompt 組裝、敏感資訊遮蔽。

## Phase 1B：OpenRouter Client 與成本追蹤

- [ ] 完成 OpenRouter models API 查詢與 `woolf models --pricing` 顯示。（M-14）
- [ ] 完成 chat completions SSE 串流解析，事件需包含 delta、done、error、usage。（M-06、M-14）
- [ ] 實作 timeout、rate limit、HTTP/API 錯誤分類與可重試錯誤處理。（M-14）
- [ ] 完成模型定價快取與 token/cost 累計。（M-17）
- [ ] 確保成本估算與 OpenRouter Dashboard 誤差目標小於 5%。（M-17）
- [ ] 測試：以 mock HTTP server 覆蓋 models、stream、錯誤事件、usage 與成本計算。

## Phase 1C：Session Schema、儲存與 CLI 管理

- [ ] 完成 Session JSON schema 與 spec 欄位對齊，包含 source、agents_config、rounds、responses、interventions、totals。（M-11）
- [ ] 實作 session ID 產生、前綴匹配與 `woolf list` 數字索引解析。（M-12、M-13）
- [ ] 完成每輪結束後自動儲存，並確保一般大小 session 儲存延遲小於 1 秒。（M-11）
- [ ] 完成 `woolf list --limit --since --status` 表格輸出。（M-13）
- [ ] 完成 `woolf show <session-id>` 唯讀檢視資料來源。（M-13）
- [ ] 完成 `woolf resume <session-id>` 從 paused/active session 載入狀態。（M-12）
- [ ] 完成 `woolf delete <session-id> [--force]`，未加 `--force` 時需確認。（M-13）
- [ ] 測試：session 儲存/讀取、ID 解析、list filter、resume 狀態、防止損壞 JSON 造成未處理崩潰。

## Phase 1D：稿件 Ingestion

- [ ] 完成 `.md` 讀取，確保中英混排與 UTF-8 不亂碼。（M-08）
- [ ] 完成 `.txt` 讀取，確保 UTF-8 純文字可直接載入。（M-09）
- [ ] 完成 ingestion 統一介面與檔案類型偵測，錯誤訊息需指出不支援格式或讀取失敗原因。（M-08、M-09、M-10）
- [ ] 完成 Phase 1 PDF 範圍：純文字層 PDF、xref、FlateDecode、基本字型、ToUnicode/CMap、文字座標排序、Markdown 轉換。（M-10）
- [ ] 明確拒絕 Phase 1 不支援的 PDF 能力：掃描/OCR、LZWDecode、多欄佈局、表格。（M-10）
- [ ] 建立 PDF regression fixture 與預期 Markdown，比對準確率目標 ≥ 85%。（M-10）
- [ ] 測試：md/txt 讀取、PDF parser 單元測試、PDF regression、50 頁內解析小於 5 秒。

## Phase 1E：Orchestrator 流水線核心

- [ ] 實作 2 到 6 Agent 依序發言，後續 Agent 可讀取前序完整內容。（M-01）
- [ ] 實作 1 到 N 輪狀態機：idle、running、waiting_intervention、paused、completed、error。（M-01、M-12、M-16）
- [ ] 實作 stance 標籤輸出與解析：agree、disagree、extend、neutral。（M-02）
- [ ] 實作 context builder，包含稿件、session 歷史、當輪前序回應、使用者介入與 focus range。（M-01、M-16）
- [ ] 實作中斷、暫停與恢復流程，Ctrl+C 或 `/pause` 後可 `resume`。（M-12）
- [ ] 實作使用者介入：輪次之間插入追問、補充、跳過 Agent、追加檔案等資料流。（M-16）
- [ ] 測試：3 Agent × 3 輪完整執行、context 傳遞順序、stance 記錄、介入內容進入下一輪、暫停續接。

## Phase 1F：TUI 三區佈局與鍵盤操作

- [ ] 完成 Bubble Tea 主 Model 與三區佈局：討論串、輸入區、狀態列。（M-05）
- [ ] 完成討論串面板：Agent badge、stance tag、串流內容、歷史 round 分隔。（M-05、M-06）
- [ ] 完成輸入區：長文輸入至少 5000 字不崩潰不卡頓。（M-05）
- [ ] 完成狀態列：目前 round、Agent、token、費用、狀態與錯誤提示。（M-17）
- [ ] 完成純鍵盤操作與 keymap：上下瀏覽、送出、命令模式、暫停、結束、說明。（M-07）
- [ ] 完成 `/start`、`/next`、`/pause`、`/end`、`/focus`、`/add-file`、`/skip`、`/summarize` 的 Phase 1 行為或明確提示未啟用。（M-16）
- [ ] 確保串流渲染延遲目標小於 100ms，串流中 TUI FPS 目標 ≥ 30。（M-06）
- [ ] 測試：TUI update model、keymap、命令解析、長文輸入、串流事件轉畫面狀態。

## Phase 1G：CLI/TUI 整合流程

- [ ] 完成 `woolf init` 互動式設定建立。
- [ ] 完成 `woolf start --draft FILE --preset NAME --agents a,b,c --rounds N` 建立 session 並進入 TUI。（M-01、M-04、M-05）
- [ ] 完成無 draft 的自由輸入流程。（M-05、M-16）
- [ ] 完成 `woolf resume` 載入 paused session 並回到 TUI。（M-12）
- [ ] 完成 `woolf show` 唯讀模式，避免誤修改 session。（M-13）
- [ ] 完成 `woolf fork <session-id> --draft FILE --title TITLE` 基礎能力，若列為 Should 可在 Phase 1H 補齊。（S-05）
- [ ] 端對端測試：不需真實 API key 的 mock OpenRouter 完整 start → 3 Agent → 儲存 → resume → export flow。

## Phase 1H：匯出、Should Have 與品質收斂

- [ ] 完成 Markdown 匯出格式：原稿、各 round、Agent 標識、stance、介入內容、統計。（M-18）
- [ ] 完成 `woolf export <session-id> --format md --output PATH`。（M-18）
- [ ] 補齊 TUI session browser，可互動式瀏覽歷史 sessions。（S-01）
- [ ] 補齊 context 摘要壓縮策略，超過 token 閾值時自動摘要舊輪次。（S-02）
- [ ] 補齊預算上限警告，達到閾值時在 TUI 顯示。（S-03）
- [ ] 補齊 `/focus` 指定段落後下一輪只針對該範圍。（S-04）
- [ ] 補齊 session fork 完整流程與獨立執行驗證。（S-05）
- [ ] 評估 PDF 匯出是否納入 Phase 1；若納入，需完成 `--format pdf`，否則保持明確未支援訊息。（S-06）
- [ ] 評估輸入區 Markdown 語法高亮是否納入 Phase 1 收斂。（S-07）

## Phase 1I：驗收、CI 與發佈準備

- [ ] 單元測試覆蓋率達到 spec 目標 ≥ 70%。
- [ ] 關鍵整合路徑測試達到 100%：CLI → Orchestrator → mocked OpenRouter → Session → Export。
- [ ] CI 執行 `go mod download`、`go test ./...`、`go vet ./...`，並涵蓋 Ubuntu 與 Windows。
- [ ] 補齊 PDF regression job，固定測試集輸出差異超過閾值即失敗。
- [ ] 驗證 Phase 1 驗收矩陣：流水線、TUI、檔案、Session、匯出、成本、設定。
- [ ] 驗證效能指標：冷啟動 < 500ms、儲存 < 1 秒、PDF 解析 < 5 秒、一般 session 記憶體 < 200MB。
- [ ] 檢查安全性：不硬編 API key、日誌遮蔽敏感值、設定權限與錯誤輸出不洩漏 token。
- [ ] 準備 v1.0.0 發佈條件：Phase 1 功能完成、Session 格式穩定、README/使用文件同步更新。

## Phase 2：後續路線圖

- [ ] P2-1：Summarizer Agent，自動產出結構化總結。
- [ ] P2-2：沙龍記憶，讓 Agent 跨 Session 累積對作者的觀察。
- [ ] P2-3：辯論模式，支援 Agent 正反對戰與裁判。
- [ ] P2-4：Micro-session，針對段落啟動小型審議。
- [ ] P2-5：Session diff，比較同稿件兩次審議結果。
- [ ] P2-6：分歧度指標，視覺化共識與爭議程度。
- [ ] P2-7：重播模式，逐字串流重播歷史 session。
- [ ] P2-8：PDF 多欄偵測。
- [ ] P2-9：PDF 匯出，優先評估 Typst 方案。
- [ ] P2-10：本地模型支援，擴展 OpenRouter Client interface 以支援 Ollama。
- [ ] P2-11：插件系統，支援第三方 Agent 角色定義分享。

## 持續性守則

- [ ] 每次功能實作都需先確認責任歸屬，不把 CLI/TUI、商業流程、API 呼叫與持久化混在同一層。
- [ ] 新增功能時同步補測試；跨模組行為至少補整合測試。
- [ ] 不因短期方便改動 Session JSON 或 CLI 輸出格式；若必須改，需先評估相容性。
- [ ] 所有外部 API 行為都應可 mock，避免測試依賴真實 OpenRouter API key。
- [ ] PDF parser 的 Phase 1 範圍要守住，不把 OCR、多欄、表格等高風險功能混入核心交付。
