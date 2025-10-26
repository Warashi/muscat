# muscat expose-port ガイド

`muscat expose-port` はローカルで動いている TCP サービスを、SSH で接続している
リモートマシンへ公開するための逆方向ポートフォワード機能です。単一の双方向
ストリーム上で複数接続を多重化し、 `--auto` オプションで新たに開いたポートを
自動的に公開できます。

```
+------------------+    SSH unix socket forward    +------------------+
| Local workstation | ============================> | Remote host       |
|  muscat expose    | <============================ |  muscat server    |
+------------------+                                +------------------+
         |                                                      ^
         | 1. ローカルサービスが 127.0.0.1:3000 で待ち受け        |
         |                                                      |
         v                                                      |
    Remote user から remote.host:3000 へアクセス <---------------+
```

## 手動公開フロー

1. SSH 設定で muscat の UNIX ソケットを転送する（例は `README.md` を参照）
2. リモート側で `muscat server` を起動しておく
3. ローカルでサービスを listen した状態で次を実行

```sh
$ muscat expose-port 3000
2025/10/26 07:12:04 exposeport: local port 3000 exposed on remote port 3000
```

複数ポートを同時に公開したい場合は `local[:remote]` 形式を並べます。

```sh
$ muscat expose-port 3000 5432:15432
```

指定したリモートポートが既に占有されている場合はエラーになります。
`--remote-policy next-free` を付けると競合時に空きポートへ自動で切り替わります。

## 自動公開 (--auto)

`--auto` を付けるとローカルマシンで listen 状態になった TCP ポートを監視し、
新規ポートを順次公開します。手動指定したポートと併用すると、明示的に指定した
ポートは確実に公開しつつ、監視対象も自動で追加されます。

```sh
$ muscat expose-port --auto \
    --auto-interval 3s \
    --auto-exclude 22 \
    --auto-exclude-process ssh \
    8080
2025/10/26 07:20:01 exposeport: local port 8080 exposed on remote port 8080
2025/10/26 07:20:04 exposeport: local port 5173 exposed on remote port 5173
```

### フィルタリング
- `--auto-exclude <port>` で特定ポートを除外
- `--auto-exclude-process <name>` でプロセス名を指定して除外
- 監視間隔は `--auto-interval` で制御（デフォルト 5s）

### リモートポート競合ポリシー
- `fail` (デフォルト): 競合時にエラーで停止
- `next-free`: 競合した場合はリモート側で空きポートを再取得
- `skip`: 競合したポートはスキップし、エラーを通知

## セキュリティの注意点
- `--public`（または `--bind-address 0.0.0.0`）を指定するとリモートマシン上で
  全てのインターフェースへ公開されます。ファイアウォールや認証設定を必ず確認
 してください。
- 監視には TCP ポート一覧を取得するための OS 権限が必要です。権限が足りない場合、
  監視対象から外れたポートは公開されません。
- 現在は TCP のみサポートしています。UDP や IPv6 は対象外です。

## 既存 `port-forward` との違い
- `muscat port-forward`: リモート -> ローカル の従来のポートフォワード
- `muscat expose-port`: ローカル -> リモート の逆方向ポート公開

互換性維持のため双方が同時に利用できます。

## トラブルシューティング
- `listen ... bind: address already in use`  
  既にリモートポートが使用中です。`--remote-policy next-free` を検討してください。
- `failed to enumerate local listeners` や権限エラーのログ  
  監視対象のソケット情報を取得できません。管理者権限または代替手段をご検討ください。
- `context canceled` のメッセージと共に終了する  
  SIGINT/SIGTERM を受けて停止しました。CLI は自動的に公開済み接続をクリーンアップします。

より詳しいテストや実装背景は `client/exposeport` と `server/exposeport` のパッケージを
参照してください。
