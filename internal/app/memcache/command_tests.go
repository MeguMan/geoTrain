package memcache

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"testing"
)

func TestLRU_Set(t *testing.T) {
	l := TestLru(t)
	l.Set("firstKey", "newValue", 0)
	newVal, _ := l.Get("firstKey")
	assert.Equal(t, "newValue", newVal)
	assert.Equal(t, l.queue.Front(), l.items["firstKey"])

	for i:= 0; i < l.capacity; i++ {
		l.Set(strconv.Itoa(i), "value", 0)
	}
	newVal, _ = l.Get("firstKey")
	assert.Equal(t, nil, newVal)
}

func TestLRU_Hset(t *testing.T) {

}

func TestLRU_Get(t *testing.T) {
	l := TestLru(t)
	val, err := l.Get("firstKey")
	assert.NoError(t, err)
	assert.NotNil(t, val)
	val1, err1 := l.Get("nonexistentKey")
	assert.Error(t, err1)
	assert.Nil(t, val1)
}

func TestLRU_GetAllKeys(t *testing.T) {
	l := TestLru(t)
	assert.Equal(t, []string{"firstKey"}, l.GetAllKeys())
}

func TestLRU_Save(t *testing.T) {
	l := TestLru(t)
	assert.NoError(t, l.Save())
	assert.FileExists(t, "data.txt")
	os.Remove("data.txt")
}

func TestLRU_Delete(t *testing.T) {
	l := TestLru(t)
	assert.NoError(t, l.Delete("firstKey"))
	assert.Error(t, l.Delete("firstKey"))
}

func TestLRU_CheckPassword(t *testing.T) {
	testCases := []struct {
		name     string
		l          *LRU
		password string
		isValid    bool
	}{
		{
			name: "valid password",
			l: TestLru(t),
			password: "strongpassword",
			isValid: true,
		},
		{
			name: "invalid password",
			l: TestLru(t),
			password: "wrongpassword",
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.True(t, tc.l.CheckPassword(tc.password))
			} else {
				assert.False(t, tc.l.CheckPassword(tc.password))
			}
		})
	}
}

func TestLRU_deleteAfterExpiration(t *testing.T) {
	l := TestLru(t)
	t.set
}