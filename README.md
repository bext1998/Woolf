# Woolf

Woolf 是一套面向文字創作者與內容工作者的多模型 AI 審議 CLI/TUI。它透過 OpenRouter API 調度多個 AI Agent，讓不同角色依序閱讀稿件、提出觀點、互相回應與辯駁，協助使用者取得更立體的創作回饋。

名稱取自 Virginia Woolf 與布魯姆斯伯里文學沙龍的傳統：一群作家、藝術家定期聚會，互相批評、激辯與啟發。WoolfCLI 的目標，是把這種多視角審議帶進終端機。

## 核心概念

單一 AI 通常只會給出一種聲音。WoolfCLI 讓使用者同時聽見嚴格編輯、一般讀者、結構分析師、支持者與挑戰者等不同視角，最後由創作者自己做決定。

預期流程：

```bash
woolf start --draft chapter3.md --preset editorial
```

啟動後，Woolf 會載入稿件與 Agent preset，進入 TUI 介面，並依照設定的輪次執行多模型審議流水線。

## 主要功能規劃

- 多 Agent 流水線：2 到 6 位 Agent 依序發言，後續 Agent 可參考前序完整內容。
- 立場標籤：Agent 可對前序觀點標記 `agree`、`disagree`、`extend` 或 `neutral`。
- 角色系統：內建多種 Agent 角色模板，並支援 YAML 自訂。
- Preset 組合：針對編輯、腦力激盪、讀者視角等情境快速啟動。
- TUI 操作：以討論串、輸入區與狀態列組成主要介面，支援鍵盤操作。
- 串流輸出：透過 OpenRouter SSE 串流即時顯示模型回應。
- 檔案載入：規劃支援 Markdown、純文字與具文字層的 PDF。
- Session 管理：自動儲存、續接、瀏覽與匯出審議紀錄。
- 成本追蹤：即時顯示 token 用量與預估費用。

## 技術方向

WoolfCLI 預計以 Go 實作，核心模組包含：

- `cmd/woolf`：CLI 入口。
- `internal/cli`：子指令與命令列參數。
- `internal/tui`：Bubble Tea TUI 介面。
- `internal/orchestrator`：Agent 流水線排程與狀態機。
- `internal/agents`：角色、prompt 與 preset 管理。
- `internal/openrouter`：OpenRouter API client、模型清單與串流處理。
- `internal/ingestion`：稿件讀取與格式轉換。
- `internal/session`：Session 持久化、續接與歷史管理。
- `internal/export`：Markdown 與後續 PDF 匯出。
- `internal/cost`：token 與費用估算。
- `internal/config`：TOML 設定檔與環境變數載入。

暫定主要依賴：

- Cobra：CLI 指令框架。
- Bubble Tea、Lip Gloss、Bubbles：TUI 介面與元件。
- OpenRouter Chat Completions API：多模型調度。

## 目標平台

- macOS
- Linux
- Windows Terminal
- WSL

## 安全原則

- 不硬編 API key、token 或任何敏感資料。
- OpenRouter API key 應存放於本機設定或環境變數。
- TUI 與日誌中顯示 API key 時必須遮蔽。
- 稿件會透過 HTTPS 傳送至 OpenRouter，使用者應自行評估敏感內容風險。

## 專案狀態

目前專案仍在早期規劃與實作前階段。公開 README 僅保留產品方向、架構概念與功能範圍；完整內部規格文件不納入 Git 追蹤。

## License

授權方式尚未決定。
