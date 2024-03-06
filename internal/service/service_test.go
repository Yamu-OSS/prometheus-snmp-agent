package service

import (
	"testing"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/stretchr/testify/assert"
)

func TestContain(t *testing.T) {
	test := []struct {
		s    labels.Labels
		v    labels.Labels
		want bool
	}{
		{
			s: labels.Labels{
				{
					Name:  "test1",
					Value: "value1",
				},
			},
			v: labels.Labels{
				{
					Name:  "test1",
					Value: "value1",
				},
			},
			want: true,
		},
		{
			s: labels.Labels{
				{
					Name:  "test1",
					Value: "value1",
				},
				{
					Name:  "test2",
					Value: "value2",
				},
			},
			v: labels.Labels{
				{
					Name:  "test1",
					Value: "value1",
				},
			},
			want: true,
		},
		{
			s: labels.Labels{
				{
					Name:  "test2",
					Value: "value2",
				},
				{
					Name:  "test1",
					Value: "value1",
				},
			},
			v: labels.Labels{
				{
					Name:  "test1",
					Value: "value1",
				},
			},
			want: true,
		},
		{
			s: labels.Labels{
				{
					Name:  "test1",
					Value: "value1",
				},
			},
			v: labels.Labels{
				{
					Name:  "test1",
					Value: "value1",
				},
				{
					Name:  "test2",
					Value: "value2",
				},
			},
			want: false,
		},
		{
			s: labels.Labels{
				{
					Name:  "test1",
					Value: "value1",
				},
			},
			v: labels.Labels{
				{
					Name:  "test2",
					Value: "value2",
				},
			},
			want: false,
		},
		{
			s: labels.Labels{
				{
					Name:  "test1",
					Value: "value1",
				},
			},
			v:    labels.Labels{},
			want: true,
		},
		{
			s: labels.Labels{},
			v: labels.Labels{
				{
					Name:  "test1",
					Value: "value1",
				},
			},
			want: false,
		},
		{
			s: labels.Labels{
				{
					Name:  "test1",
					Value: "value3",
				},
			},
			v: labels.Labels{
				{
					Name:  "test1",
					Value: "value1",
				},
			},
			want: false,
		},
		{
			s: labels.Labels{
				{
					Name:  "test1",
					Value: "value1",
				},
				{
					Name:  "test2",
					Value: "value2",
				},
				{
					Name:  "test3",
					Value: "value3",
				},
			},
			v: labels.Labels{
				{
					Name:  "test1",
					Value: "value1",
				},
				{
					Name:  "test2",
					Value: "value2",
				},
			},
			want: true,
		},
	}

	for _, tv := range test {
		get := Contains(tv.s, tv.v)
		assert.Equal(t, tv.want, get)
	}
}
