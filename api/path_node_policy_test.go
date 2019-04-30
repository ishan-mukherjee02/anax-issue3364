// +build unit

package api

import (
	"flag"
	"github.com/open-horizon/anax/externalpolicy"
	_ "github.com/open-horizon/anax/externalpolicy/text_language"
	"testing"
)

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("v", "7")
	// no need to parse flags, that's done by test framework
}

// Verify that FindNodePolicyForOutput works when there is no node policy defined yet.
func Test_FindNPForOutput0(t *testing.T) {

	dir, db, err := utsetup()
	if err != nil {
		t.Error(err)
	}
	defer cleanTestDir(dir)

	if np, err := FindNodePolicyForOutput(db); err != nil {
		t.Errorf("failed to find node policy in db, error %v", err)
	} else if len(np.Properties) != 0 {
		t.Errorf("incorrect node policy, there should be no properties defined, found: %v", *np)
	}

}

// Verify that a Node Policy Object can be created and saved the first time.
func Test_SaveNodePolicy1(t *testing.T) {

	dir, db, err := utsetup()
	if err != nil {
		t.Error(err)
	}
	defer cleanTestDir(dir)

	var myError error
	errorhandler := GetPassThroughErrorHandler(&myError)

	propName := "prop1"
	propList := new(externalpolicy.PropertyList)
	propList.Add_Property(externalpolicy.Property_Factory(propName, "val1"))

	extNodePolicy := &externalpolicy.ExternalPolicy{
		Properties:  *propList,
		Constraints: []string{`prop3 == "some value"`},
	}

	errHandled, np, msgs := UpdateNodePolicy(extNodePolicy, errorhandler, db, getBasicConfig())

	if errHandled {
		t.Errorf("Unexpected error handled: %v", myError)
	} else if myError != nil && len(myError.Error()) != 0 {
		t.Errorf("myError set unexpectedly (%T) %v", myError, myError)
	} else if np == nil {
		t.Errorf("no node policy returned")
	} else if fnp, err := FindNodePolicyForOutput(db); err != nil {
		t.Errorf("failed to find node policy in db, error %v", err)
	} else if len(fnp.Properties) != 1 {
		t.Errorf("incorrect node policy, there should be 1 property defined, found: %v", *fnp)
	} else if fnp.Properties[0].Name != propName {
		t.Errorf("expected property %v, but received %v", propName, fnp.Properties[0].Name)
	} else if len(msgs) != 1 {
		t.Errorf("there should be 1 message, returned %v", len(msgs))
	}

}

// Verify that a Node Policy Object can be created and saved, and then updated.
func Test_UpdateNodePolicy1(t *testing.T) {

	dir, db, err := utsetup()
	if err != nil {
		t.Error(err)
	}
	defer cleanTestDir(dir)

	var myError error
	errorhandler := GetPassThroughErrorHandler(&myError)

	propName := "prop1"
	propList := new(externalpolicy.PropertyList)
	propList.Add_Property(externalpolicy.Property_Factory(propName, "val1"))

	extNodePolicy := &externalpolicy.ExternalPolicy{
		Properties:  *propList,
		Constraints: []string{`prop3 == "some value"`},
	}

	errHandled, np, msgs := UpdateNodePolicy(extNodePolicy, errorhandler, db, getBasicConfig())

	if errHandled {
		t.Errorf("Unexpected error handled: %v", myError)
	} else if myError != nil && len(myError.Error()) != 0 {
		t.Errorf("myError set unexpectedly (%T) %v", myError, myError)
	} else if np == nil {
		t.Errorf("no node policy returned")
	} else if fnp, err := FindNodePolicyForOutput(db); err != nil {
		t.Errorf("failed to find node policy in db, error %v", err)
	} else if len(fnp.Properties) != 1 {
		t.Errorf("incorrect node policy, there should be 1 property defined, found: %v", *fnp)
	} else if fnp.Properties[0].Name != propName {
		t.Errorf("expected property %v, but received %v", propName, fnp.Properties[0].Name)
	} else if len(msgs) != 1 {
		t.Errorf("there should be 1 message, returned %v", len(msgs))
	}

	// Now change the property specified in the policy.
	propName = "prop2"
	propList = new(externalpolicy.PropertyList)
	propList.Add_Property(externalpolicy.Property_Factory(propName, "val2"))

	extNodePolicy.Properties = *propList

	errHandled, np, msgs = UpdateNodePolicy(extNodePolicy, errorhandler, db, getBasicConfig())

	if errHandled {
		t.Errorf("Unexpected error handled: %v", myError)
	} else if myError != nil && len(myError.Error()) != 0 {
		t.Errorf("myError set unexpectedly (%T) %v", myError, myError)
	} else if np == nil {
		t.Errorf("no node policy returned")
	} else if fnp, err := FindNodePolicyForOutput(db); err != nil {
		t.Errorf("failed to find node policy in db, error %v", err)
	} else if len(fnp.Properties) != 1 {
		t.Errorf("incorrect node policy, there should be 1 property defined, found: %v", *fnp)
	} else if fnp.Properties[0].Name != propName {
		t.Errorf("expected property %v, but received %v", propName, fnp.Properties[0].Name)
	} else if len(msgs) != 1 {
		t.Errorf("there should be 1 message, returned %v", len(msgs))
	}

}

// Verify that a Node Policy Object can be created and then deleted.
func Test_DeleteNodePolicy1(t *testing.T) {

	dir, db, err := utsetup()
	if err != nil {
		t.Error(err)
	}
	defer cleanTestDir(dir)

	var myError error
	errorhandler := GetPassThroughErrorHandler(&myError)

	propName := "prop1"
	propList := new(externalpolicy.PropertyList)
	propList.Add_Property(externalpolicy.Property_Factory(propName, "val1"))

	extNodePolicy := &externalpolicy.ExternalPolicy{
		Properties:  *propList,
		Constraints: []string{`prop3 == "some value"`},
	}

	errHandled, np, msgs := UpdateNodePolicy(extNodePolicy, errorhandler, db, getBasicConfig())

	if errHandled {
		t.Errorf("Unexpected error handled: %v", myError)
	} else if myError != nil && len(myError.Error()) != 0 {
		t.Errorf("myError set unexpectedly (%T) %v", myError, myError)
	} else if np == nil {
		t.Errorf("no node policy returned")
	} else if fnp, err := FindNodePolicyForOutput(db); err != nil {
		t.Errorf("failed to find node policy in db, error %v", err)
	} else if len(fnp.Properties) != 1 {
		t.Errorf("incorrect node policy, there should be 1 property defined, found: %v", *fnp)
	} else if fnp.Properties[0].Name != propName {
		t.Errorf("expected property %v, but received %v", propName, fnp.Properties[0].Name)
	} else if len(msgs) != 1 {
		t.Errorf("there should be 1 message, returned %v", len(msgs))
	}

	// Now delete the object.

	errHandled, msgs = DeleteNodePolicy(errorhandler, db)

	if errHandled {
		t.Errorf("Unexpected error handled: %v", myError)
	} else if myError != nil && len(myError.Error()) != 0 {
		t.Errorf("myError set unexpectedly (%T) %v", myError, myError)
	} else if np == nil {
		t.Errorf("no node policy returned")
	} else if fnp, err := FindNodePolicyForOutput(db); err != nil {
		t.Errorf("failed to find node policy in db, error %v", err)
	} else if len(fnp.Properties) != 0 {
		t.Errorf("incorrect node policy, there should be 0 properties defined, found: %v", *fnp)
	} else if len(msgs) != 1 {
		t.Errorf("there should be 1 message, returned %v", len(msgs))
	}

}
