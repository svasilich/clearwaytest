package asset

import "context"

// DBAssetWriter write asset data to data base.
type DBAssetWriter interface {
	WriteAsset(ctx context.Context, asset string, uid int64, data []byte) error
}

// DBAssetReader read asset data from data base.
type DBAssetReader interface {
	ReadAsset(ctx context.Context, asset string, uid int64) ([]byte, error)
}
