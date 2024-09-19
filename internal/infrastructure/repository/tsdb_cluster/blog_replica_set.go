package tsdb_cluster

import (
	"info/internal/domain/blog"
	"info/internal/infrastructure/repository/tsdb"
)

type BlogReplicaSet struct {
	*ReplicaSet
}

var _ blog.ReplicaSet = (*BlogReplicaSet)(nil)

func NewBlogReplicaSet(replicaSet *ReplicaSet) *BlogReplicaSet {
	return &BlogReplicaSet{
		ReplicaSet: replicaSet,
	}
}

func (c *BlogReplicaSet) WriteRepo() blog.WriteRepository {
	return tsdb.NewBlogRepository(c.ReplicaSet.WriteRepo())
}

func (c *BlogReplicaSet) ReadRepo() blog.ReadRepository {
	return tsdb.NewBlogRepository(c.ReplicaSet.ReadRepo())
}
