package lduser

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/launchdarkly/go-sdk-common.v2/ldvalue"
)

type userStringPropertyDesc struct {
	name   string
	getter func(User) ldvalue.OptionalString
	setter func(UserBuilder, string) UserBuilderCanMakeAttributePrivate
}

var userSecondaryKeyProperty = userStringPropertyDesc{
	"secondary",
	User.GetSecondaryKey,
	UserBuilder.Secondary,
}
var userIPProperty = userStringPropertyDesc{
	"ip",
	User.GetIP,
	UserBuilder.IP,
}
var userCountryProperty = userStringPropertyDesc{
	"country",
	User.GetCountry,
	UserBuilder.Country,
}
var emailProperty = userStringPropertyDesc{
	"country",
	User.GetCountry,
	UserBuilder.Country,
}
var userFirstNameProperty = userStringPropertyDesc{
	"firstName",
	User.GetFirstName,
	UserBuilder.FirstName,
}
var userLastNameProperty = userStringPropertyDesc{
	"lastName",
	User.GetLastName,
	UserBuilder.LastName,
}
var userAvatarProperty = userStringPropertyDesc{
	"avatar",
	User.GetAvatar,
	UserBuilder.Avatar,
}
var userNameProperty = userStringPropertyDesc{
	"name",
	User.GetName,
	UserBuilder.Name,
}

var allUserStringProperties = []userStringPropertyDesc{
	userSecondaryKeyProperty,
	userIPProperty,
	userCountryProperty,
	userFirstNameProperty,
	userLastNameProperty,
	userAvatarProperty,
	userNameProperty,
}

func (p userStringPropertyDesc) assertNotSet(t *testing.T, user User) {
	assert.Equal(t, ldvalue.OptionalString{}, p.getter(user), "should not have had a value for %s", p.name)
}

func assertStringPropertiesNotSet(t *testing.T, user User) {
	for _, p := range allUserStringProperties {
		p.assertNotSet(t, user)
	}
}

func TestNewUser(t *testing.T) {
	user := NewUser("some-key")

	assert.Equal(t, "some-key", user.GetKey())

	k, ok := user.valueOf("key")
	assert.True(t, ok)
	assert.Equal(t, ldvalue.String("some-key"), k)

	for _, p := range allUserStringProperties {
		p.assertNotSet(t, user)
	}

	anon, found := user.GetAnonymousOptional()
	assert.False(t, anon)
	assert.False(t, found)

	assert.Nil(t, user.Custom)
	assert.Nil(t, user.PrivateAttributeNames)
}

func TestNewAnonymousUser(t *testing.T) {
	user := NewAnonymousUser("some-key")

	assert.Equal(t, "some-key", user.GetKey())

	k, ok := user.valueOf("key")
	assert.True(t, ok)
	assert.Equal(t, ldvalue.String("some-key"), k)

	for _, p := range allUserStringProperties {
		p.assertNotSet(t, user)
	}

	anon, found := user.GetAnonymousOptional()
	assert.True(t, anon)
	assert.True(t, found)

	assert.Nil(t, user.Custom)
	assert.Nil(t, user.PrivateAttributeNames)
}

func TestUserBuilderSetsOnlyKeyByDefault(t *testing.T) {
	user := NewUserBuilder("some-key").Build()

	assert.Equal(t, "some-key", user.GetKey())

	k, _ := user.valueOf("key")
	assert.Equal(t, ldvalue.String("some-key"), k)

	for _, p := range allUserStringProperties {
		p.assertNotSet(t, user)
	}

	anon, found := user.GetAnonymousOptional()
	assert.False(t, anon)
	assert.False(t, found)

	assert.Nil(t, user.Custom)
	assert.Nil(t, user.PrivateAttributeNames)
}

func TestUserBuilderCanSetStringAttributes(t *testing.T) {
	for _, p := range allUserStringProperties {
		t.Run(p.name, func(t *testing.T) {
			builder := NewUserBuilder("some-key")
			p.setter(builder, "value")
			user := builder.Build()

			assert.Equal(t, "some-key", user.GetKey())

			for _, p1 := range allUserStringProperties {
				if p1.name == p.name {
					assert.Equal(t, ldvalue.NewOptionalString("value"), p.getter(user), p.name)
					v, ok := user.valueOf(p.name)
					if p.name == "secondary" {
						// this attribute is special in that it *can't* be used in evaluations
						assert.False(t, ok, p.name)
						assert.Equal(t, ldvalue.Null(), v, p.name)
					} else {
						assert.True(t, ok, p.name)
						assert.Equal(t, ldvalue.String("value"), v)
					}
				} else {
					p1.assertNotSet(t, user)
					v, ok := user.valueOf(p1.name)
					assert.False(t, ok, p1.name)
					assert.Equal(t, ldvalue.Null(), v, p1.name)
				}
			}
		})
	}
}

func TestUserBuilderCanSetAnonymous(t *testing.T) {
	user0 := NewUserBuilder("some-key").Build()
	assert.False(t, user0.GetAnonymous())
	value, ok := user0.GetAnonymousOptional()
	assert.False(t, ok)
	assert.False(t, value)

	user1 := NewUserBuilder("some-key").Anonymous(true).Build()
	assert.True(t, user1.GetAnonymous())
	value, ok = user1.GetAnonymousOptional()
	assert.True(t, ok)
	assert.True(t, value)

	user2 := NewUserBuilder("some-key").Anonymous(false).Build()
	assert.False(t, user2.GetAnonymous())
	value, ok = user2.GetAnonymousOptional()
	assert.True(t, ok)
	assert.False(t, value)
	assert.NotNil(t, user2.Anonymous)
}

func TestUserBuilderCanSetPrivateStringAttributes(t *testing.T) {
	for _, p := range allUserStringProperties {
		t.Run(p.name, func(t *testing.T) {
			builder := NewUserBuilder("some-key")
			p.setter(builder, "value").AsPrivateAttribute()
			user := builder.Build()

			assert.Equal(t, "some-key", user.GetKey())

			for _, p1 := range allUserStringProperties {
				if p1.name == p.name {
					assert.Equal(t, ldvalue.NewOptionalString("value"), p.getter(user))
				} else {
					p1.assertNotSet(t, user)
				}
			}

			assert.Nil(t, user.Custom)
			assert.Equal(t, []string{p.name}, user.PrivateAttributeNames)
		})
	}
}

func TestUserBuilderCanMakeAttributeNonPrivate(t *testing.T) {
	builder := NewUserBuilder("some-key")
	builder.Country("us").AsNonPrivateAttribute()
	builder.Email("e").AsPrivateAttribute()
	builder.Name("n").AsPrivateAttribute()
	builder.Email("f").AsNonPrivateAttribute()
	user := builder.Build()
	assert.Equal(t, "f", user.GetEmail().StringValue())
	assert.Equal(t, []string{"name"}, user.PrivateAttributeNames)
}

func TestUserBuilderCanSetCustomAttributes(t *testing.T) {
	user := NewUserBuilder("some-key").Custom("first", ldvalue.Int(1)).Custom("second", ldvalue.String("two")).Build()

	value, ok := user.GetCustom("first")
	assert.True(t, ok)
	assert.Equal(t, 1, value.IntValue())

	value, ok = user.GetCustom("second")
	assert.True(t, ok)
	assert.Equal(t, "two", value.StringValue())

	value, ok = user.GetCustom("no")
	assert.False(t, ok)
	assert.Equal(t, ldvalue.Null(), value)

	keys := user.GetCustomKeys()
	sort.Strings(keys)
	assert.Equal(t, []string{"first", "second"}, keys)

	assert.Nil(t, user.PrivateAttributeNames)
}

func TestUserWithNoCustomAttributes(t *testing.T) {
	user := NewUser("some-key")

	assert.Nil(t, user.Custom)

	value, ok := user.GetCustom("attr")
	assert.False(t, ok)
	assert.Equal(t, ldvalue.Null(), value)

	assert.Nil(t, user.GetCustomKeys())
}

func TestUserBuilderCanSetPrivateCustomAttributes(t *testing.T) {
	user := NewUserBuilder("some-key").Custom("first", ldvalue.Int(1)).AsPrivateAttribute().
		Custom("second", ldvalue.String("two")).Build()

	value, ok := user.GetCustom("first")
	assert.True(t, ok)
	assert.Equal(t, 1, value.IntValue())

	value, ok = user.GetCustom("second")
	assert.True(t, ok)
	assert.Equal(t, "two", value.StringValue())

	value, ok = user.GetCustom("no")
	assert.False(t, ok)
	assert.Equal(t, ldvalue.Null(), value)

	keys := user.GetCustomKeys()
	sort.Strings(keys)
	assert.Equal(t, []string{"first", "second"}, keys)

	assert.NotNil(t, user.PrivateAttributeNames)
	assert.Equal(t, []string{"first"}, user.PrivateAttributeNames)
}

func TestUserBuilderCanCopyFromExistingUserWithOnlyKey(t *testing.T) {
	user0 := NewUser("some-key")
	user1 := NewUserBuilderFromUser(user0).Build()

	assert.Equal(t, "some-key", user1.GetKey())

	for _, p := range allUserStringProperties {
		p.assertNotSet(t, user1)
	}
	assert.Nil(t, user1.Custom)
	assert.Nil(t, user1.PrivateAttributeNames)
}

func TestUserBuilderCanCopyFromExistingUserWithAllAttributes(t *testing.T) {
	user0 := newUserBuilderWithAllPropertiesSet("some-key").Build()
	user1 := NewUserBuilderFromUser(user0).Build()
	assert.Equal(t, user0, user1)
}

func TestUserEqualsComparesAllAttributes(t *testing.T) {
	shouldNotEqual := func(a User, b User) {
		assert.False(t, b.Equal(a), "%s should not equal %s", b, a)
	}

	user0 := NewUser("some-key")
	assert.True(t, user0.Equal(user0), "%s should equal itself", user0)

	user1 := newUserBuilderWithAllPropertiesSet("some-key").Build()
	assert.True(t, user1.Equal(user1), "%s should equal itself", user1)
	user2 := NewUserBuilderFromUser(user1).Build()
	assert.True(t, user2.Equal(user1), "%s should equal %s", user2, user1)

	for i, p := range allUserStringProperties {
		builder3 := NewUserBuilderFromUser(user1)
		p.setter(builder3, "different-value")
		user3 := builder3.Build()
		shouldNotEqual(user1, user3)

		builder4 := NewUserBuilderFromUser(user1)
		p.setter(builder4, fmt.Sprintf("value%d", i)).AsPrivateAttribute()
		user4 := builder4.Build()
		shouldNotEqual(user1, user4)
	}

	shouldNotEqual(user1, NewUserBuilderFromUser(user1).Key("other-key").Build())

	shouldNotEqual(user0, NewUserBuilderFromUser(user0).Anonymous(true).Build())
	shouldNotEqual(NewUserBuilderFromUser(user0).Anonymous(true).Build(), NewUserBuilderFromUser(user0).Anonymous(false).Build())

	shouldNotEqual(user1, NewUserBuilderFromUser(user1).Custom("thing1", ldvalue.String("value9")).Build())
	shouldNotEqual(user1, NewUserBuilderFromUser(user1).Custom("thing1", ldvalue.String("value1")).AsPrivateAttribute().Build())
}

func newUserBuilderWithAllPropertiesSet(key string) UserBuilder {
	builder := NewUserBuilder(key)
	for i, p := range allUserStringProperties {
		p.setter(builder, fmt.Sprintf("value%d", i))
	}
	builder.Anonymous(true)
	builder.Custom("thing1", ldvalue.String("value1"))
	builder.Custom("thing2", ldvalue.String("value2")).AsPrivateAttribute()
	return builder
}
