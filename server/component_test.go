package server

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecursivelyFindField_NoFieldPath(t *testing.T) {
	component := reflect.ValueOf(struct {
		IgnorePlatform bool
	}{
		IgnorePlatform: false,
	})
	fieldPath := ""

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"IgnorePlatform": "bool",
	}, field)
}

func TestRecursivelyFindField_Simple(t *testing.T) {
	component := reflect.ValueOf(struct {
		IgnorePlatform bool
	}{
		IgnorePlatform: false,
	})
	fieldPath := "IgnorePlatform"

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, false, field)
}

func TestRecursivelyFindField_Nested(t *testing.T) {
	component := reflect.ValueOf(struct {
		Platform struct {
			Ignore bool
		}
	}{
		Platform: struct {
			Ignore bool
		}{
			Ignore: true,
		},
	})
	fieldPath := "Platform.Ignore"

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, true, field)
}

func TestRecursivelyFindField_NestedSlice(t *testing.T) {
	component := reflect.ValueOf(struct {
		Platforms []struct {
			Ignore bool
		}
	}{
		Platforms: []struct {
			Ignore bool
		}{
			{
				Ignore: true,
			},
		},
	})
	fieldPath := "Platforms[0].Ignore"

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, true, field)
}

func TestRecursivelyFindField_SliceInSlice(t *testing.T) {
	component := reflect.ValueOf(struct {
		Platforms [][]struct {
			Points []float64
		}
	}{
		Platforms: [][]struct {
			Points []float64
		}{
			{
				{
					Points: []float64{float64(1), float64(2)},
				},
			},
		},
	})
	fieldPath := "Platforms[0][0].Points[0]"

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, float64(1), field)
}

func TestRecursivelyFindField_SliceInSliceOfPointers(t *testing.T) {
	component := reflect.ValueOf(struct {
		Platforms [][]*struct {
			Points []float64
		}
	}{
		Platforms: [][]*struct {
			Points []float64
		}{
			{
				{
					Points: []float64{float64(1), float64(2)},
				},
			},
		},
	})

	// TODO: does this make sense to just say "ptr" here?
	// I think it would be better to return "pointer of ..." etc.
	fieldPath := "Platforms[0]" // should return slice of ptr
	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, []interface{}{"ptr"}, field)

	fieldPath = "Platforms[0][0].Points[0]"
	field, err = GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, float64(1), field)

	fieldPath = "Platforms[0][0]" // should return the struct
	field, err = GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"Points": "slice of float64",
	}, field)

	fieldPath = "Platforms[0][0].Points" // should return the slice
	field, err = GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, []interface{}([]interface{}{"float64", "float64"}), field)

}

func TestRecursivelyFindField_NestedMap(t *testing.T) {
	component := reflect.ValueOf(struct {
		Platforms map[string]struct {
			Ignore bool
		}
	}{
		Platforms: map[string]struct {
			Ignore bool
		}{
			"test": {
				Ignore: true,
			},
		},
	})
	fieldPath := "Platforms[test].Ignore"

	field, err := GetField(component, fieldPath)
	assert.Equal(t, true, field)
	assert.Nil(t, err)
}

func TestRecursivelyFindField_InvalidField(t *testing.T) {
	component := reflect.ValueOf(struct {
		IgnorePlatform bool
	}{
		IgnorePlatform: false,
	})
	fieldPath := "InvalidField"

	field, err := GetField(component, fieldPath)
	assert.Equal(t, "invalid field access", err.Error())
	assert.Nil(t, field)
}

func TestRecursivelyFindField_InvalidIndex(t *testing.T) {
	component := reflect.ValueOf(struct {
		Platforms []struct {
			Ignore bool
		}
	}{
		Platforms: []struct {
			Ignore bool
		}{
			{
				Ignore: true,
			},
		},
	})
	fieldPath := "Platforms[1].Ignore"

	field, err := GetField(component, fieldPath)
	assert.Equal(t, "invalid slice index", err.Error())
	assert.Nil(t, field)
}

func TestRecursivelyFindField_InvalidKey(t *testing.T) {
	component := reflect.ValueOf(struct {
		Platforms map[string]struct {
			Ignore bool
		}
	}{
		Platforms: map[string]struct {
			Ignore bool
		}{
			"test": {
				Ignore: true,
			},
		},
	})
	fieldPath := "Platforms[invalid].Ignore"

	field, err := GetField(component, fieldPath)
	assert.Equal(t, "invalid map key", err.Error())
	assert.Nil(t, field)
}

func TestRecursivelyFindField_NilPointer(t *testing.T) {
	component := reflect.ValueOf(struct {
		Platform *struct {
			Ignore bool
		}
	}{
		Platform: nil,
	})
	fieldPath := "Platform.Ignore"

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Nil(t, field)
}

func TestRecursivelyFindField_StructPtr(t *testing.T) {
	component := reflect.ValueOf(struct {
		Platform *struct {
			Ignore bool
		}
	}{
		Platform: &struct {
			Ignore bool
		}{
			Ignore: true,
		},
	})
	fieldPath := "Platform.Ignore"

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, true, field)
}

func TestRecursivelyFindField_StructPtrWithSlice(t *testing.T) {
	component := reflect.ValueOf(struct {
		Platform *struct {
			Ignore bool
			Points []float64
		}
	}{
		Platform: &struct {
			Ignore bool
			Points []float64
		}{
			Ignore: true,
			Points: []float64{float64(1), float64(2)},
		},
	})
	fieldPath := "Platform.Points[0]"

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, float64(1), field)
}

func TestRecursivelyFindField_DeeplyNestedStruct(t *testing.T) {
	component := reflect.ValueOf(struct {
		Platform struct {
			Settings struct {
				Ignore bool
			}
		}
	}{
		Platform: struct {
			Settings struct {
				Ignore bool
			}
		}{
			Settings: struct {
				Ignore bool
			}{
				Ignore: true,
			},
		},
	})
	fieldPath := "Platform.Settings"

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"Ignore": "bool",
	}, field)
}

func TestRecursivelyFindField_DeeplyNestedInterface(t *testing.T) {
	component := reflect.ValueOf(struct {
		Platform struct {
			Settings interface{}
		}
	}{
		Platform: struct {
			Settings interface{}
		}{
			Settings: struct {
				Ignore bool
			}{
				Ignore: true,
			},
		},
	})
	fieldPath := "Platform.Settings"

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, map[string]interface{}{
		"Ignore": "bool",
	}, field)
}

func TestRecursivelyFindField_UnexportedField(t *testing.T) {
	component := reflect.ValueOf(struct {
		Entity struct {
			ignore bool
		}
	}{
		Entity: struct {
			ignore bool
		}{
			ignore: true,
		},
	})
	fieldPath := "Entity.ignore"

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	// TODO: does this make sense to just say "Unexported field" here?
	// The CLI should probably make a distinction between unexported fields and
	// normal string values.
	assert.Equal(t, `"Unexported field"`, field)
}

func TestRecursivelyFindField_MapAccess(t *testing.T) {
	type Whatever struct {
		Ignore bool
	}
	mapToPtr := map[string]*Whatever{
		"Ptr": {
			Ignore: true,
		},
	}
	component := reflect.ValueOf(mapToPtr)
	fieldPath := "Ptr.Ignore"

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, true, field)
}

func TestRecursivelyFindField_PtrAccess(t *testing.T) {
	type Whatever struct {
		Points []float64
	}
	ptrToWhatever := &Whatever{
		Points: []float64{float64(1), float64(2)},
	}
	component := reflect.ValueOf(ptrToWhatever)
	fieldPath := "Points[0]"

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, float64(1), field)

	fieldPath = "Points[1]"
	field, err = GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, float64(2), field)

	fieldPath = "Points[2]"
	field, err = GetField(component, fieldPath)
	assert.Equal(t, "invalid slice index", err.Error())
	assert.Nil(t, field)

	fieldPath = "Points"
	field, err = GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, []interface{}([]interface{}{"float64", "float64"}), field)
}

func TestRecursivelyFindField_MultipleIndicesAccess(t *testing.T) {
	type Whatever struct {
		Points []struct {
			Points []float64
		}
	}
	component := reflect.ValueOf(Whatever{
		Points: []struct {
			Points []float64
		}{
			{
				Points: []float64{float64(1), float64(2)},
			},
		},
	})
	fieldPath := "Points[0].Points[0]"

	field, err := GetField(component, fieldPath)
	assert.Nil(t, err)
	assert.Equal(t, float64(1), field)
}
