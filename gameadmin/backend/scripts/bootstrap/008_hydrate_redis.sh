#!/usr/bin/env bash
set -euo pipefail

APP_ID="asiagaming"
EXPECTED_COUNT="413"
MONGO_PASSWORD="$(< /root/.skspin-mongo-password)"
EXPORT_FILE="$(mktemp)"

cleanup() {
  unset MONGO_PASSWORD
  rm -f "$EXPORT_FILE"
}
trap cleanup EXIT

export MONGO_PASSWORD

mongosh --quiet \
  --host 127.0.0.1 \
  --username skspin_app \
  --password "$MONGO_PASSWORD" \
  --authenticationDatabase admin \
  --eval '
const appId = "asiagaming";
const adminDb = db.getSiblingDB("GameAdmin");

adminDb.GameConfig
  .find({ AppID: appId })
  .sort({ GameId: 1 })
  .forEach(config => {
    const cs = String(
      config.BetBase || "0.1,0.2,0.3"
    )
      .split(",")
      .map(value => Number(value))
      .filter(value => Number.isFinite(value));

    const payload = {
      reward_percent: Number(config.RewardPercent || 0),
      no_award_percent: Number(config.NoAwardPercent || 0),
      appId: appId,
      gamePatten: Number(config.GamePattern || 3),
      MaxWinPoints: Number(config.MaxWinPoints || 1000000),
      maxMultiple: Number(config.MaxMultiple || 100),
      gameId: String(config.GameId),
      cs: cs,
      StopLoss: Number(config.StopLoss || 0),
      ShowNameAndTimeOff:
        Number(config.ShowNameAndTimeOff || 0),
      ShowExitBtnOff:
        Number(config.ShowExitBtnOff || 0),
      FreeGameOff: Number(config.FreeGameOff || 0),
      DefaultCs:
        Number(config.DefaultCs || cs[0] || 0.1),
      DefaultBetLevel:
        Number(config.DefaultBetLevel || 1),
      IsProtection:
        Number(config.IsProtection || 0),
      ProtectionRotateCount:
        Number(config.ProtectionRotateCount || 0),
      ProtectionRewardPercentLess:
        Number(config.ProtectionRewardPercentLess || 0),
      BuyMinAwardPercent:
        Number(config.BuyMinAwardPercent || 0),
      RTP: Number(config.RTP || 93),
      ProfitMargin: Number(config.ProfitMargin || 3),
      CrashRate: Number(config.CrashRate || 1),
      Scale: Number(config.Scale || 3),
      OnlineUpNum: Number(config.OnlineUpNum || 0),
      OnlineDownNum: Number(config.OnlineDownNum || 0)
    };

    print(
      appId + ":" + config.GameId +
      "\t" +
      JSON.stringify(payload)
    );
  });
' > "$EXPORT_FILE"

EXPORT_COUNT="$(wc -l < "$EXPORT_FILE")"

if [ "$EXPORT_COUNT" -ne "$EXPECTED_COUNT" ]; then
  echo "ERROR: GameConfig export count is ${EXPORT_COUNT}"
  exit 1
fi

while IFS=$'\t' read -r redis_key redis_payload; do
  redis-cli SETEX \
    "$redis_key" \
    86400 \
    "$redis_payload" \
    >/dev/null
done < "$EXPORT_FILE"

HYDRATED_COUNT="$(
  redis-cli --scan \
    --pattern "${APP_ID}:*" |
  wc -l
)"

if [ "$HYDRATED_COUNT" -ne "$EXPECTED_COUNT" ]; then
  echo "ERROR: Redis hydration count is ${HYDRATED_COUNT}"
  exit 1
fi

redis-cli SET \
  "skspin:migration:008_redis_hydration" \
  "$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
  >/dev/null

echo "STEP 8 REDIS HYDRATION VALID: ${HYDRATED_COUNT}"
