package asset

import "context"

// DBAssetWriter write asset data to data base.
type DBAssetWriter interface {
	WriteAsset(ctx context.Context, asset string, uid int64, data []byte) error
}
