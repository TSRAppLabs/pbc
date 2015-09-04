package pbc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	. "stash.tsrapplabs.com/ut/jsonlint"
)

var checkPass TypeCheck

func LintPass(path string) (Warning, error) {
	file, err := os.Open(filepath.Join(path, "pass.json"))

	if err != nil {
		return []string{}, err
	}

	var content interface{}
	err = json.NewDecoder(file).Decode(&content)

	if err != nil {
		return []string{}, err
	}

	warn := CheckPass(content)

	image_warn, err := LintImages(path)

	if err != nil {
		return []string{}, err
	}

	return append(warn, image_warn...), nil

}

func CheckPass(payload interface{}) Warning {
	return checkPass(payload)
}

func isDate(val interface{}) Warning {
	return []string{}
}

var colorRegex = regexp.MustCompile("^#[0-9a-fA-F]{6}$")

func isColor(val interface{}) Warning {
	str, isStr := val.(string)

	if !isStr {
		return NewWarning("expected color as string")
	}

	if !colorRegex.MatchString(str) {
		return NewWarning("cannot detect color")
	}

	return []string{}
}

func isU16(val interface{}) Warning {
	i, isNumber := val.(float64)

	if !isNumber {
		return NewWarning("expected int")
	}

	if i > 65535 {
		return NewWarning("value exceeds limit")
	}
	return []string{}
}

func enum(base string, cases ...string) []string {
	result := []string{}

	for _, c := range cases {
		result = append(result, base+c)
	}

	return result
}

var encodings = enum("PKBarcodeFormat", "QR", "PDF417", "Aztec")
var detectorTypes = enum("PKDataDetectorType", "PhoneNumber", "Link", "Address", "CalendarEvent")
var textAlignments = enum("PKTextAlignment", "Left", "Center", "Right", "Natural")
var dateStyle = enum("PKDateStyle", "None", "Short", "Medium", "Long", "Full")
var numberStyle = enum("PKNumberStyle", "Decimal", "Percent", "Scientific", "SpellOut")

//Types
var (
	beaconType = And(Required("proximityUUID"), WhiteList("major", "minor", "proximityUUID", "relevantText"), Object(map[string]TypeCheck{
		"major":         isU16,
		"minor":         isU16,
		"proximityUUID": IsString,
		"relevantText":  IsString,
	}))

	dateStyleType = Object(map[string]TypeCheck{
		"dateStyle":       StringEnum(dateStyle...),
		"ignoresTimeZone": IsBool,
		"isRelative":      IsBool,
		"timeStyle":       StringEnum(dateStyle...),
	})

	numberStyleType = Object(map[string]TypeCheck{
		"currencyCode": IsString,
		"numberStyle":  StringEnum(numberStyle...),
	})

	fieldType = And(Required("key"), dateStyleType, numberStyleType, Object(map[string]TypeCheck{
		"attributedValue":   Either(IsString, IsNumber),
		"changeMessage":     IsString,
		"dataDetectorTypes": ArrayOf(StringEnum(detectorTypes...)),
		"key":               IsString,
		"label":             IsString,
		"textAlignment":     StringEnum(textAlignments...),
		"value":             Either(IsString, isDate, IsNumber, IsDouble),
	}), WhiteList("dateStyle", "ignoresTimeZone", "isRelative", "timeStyle", "currencyCode", "numberStyle", "attributedValue", "changeMessage", "dataDetectorTypes", "key", "label", "textAlignment", "value"))

	passStructureType = And(Object(map[string]TypeCheck{
		"auxiliaryFields": ArrayOf(fieldType),
		"backFields":      ArrayOf(fieldType),
		"headerFields":    ArrayOf(fieldType),
		"primaryFields":   ArrayOf(fieldType),
		"secondaryFields": ArrayOf(fieldType),
		//transitType (required for boarding passes), is outside of our purposes for this program currently
	}), WhiteList("auxiliaryFields", "backFields", "headerFields", "primaryFields", "secondaryFields", "transitType"))

	locationType = And(Object(map[string]TypeCheck{
		"altitude":     IsDouble,
		"latitude":     IsDouble,
		"longitude":    IsDouble,
		"relevantText": IsString,
	}), Required("latitude", "longitude"),
		WhiteList("altitude", "latitude", "longitude", "relevantText"))

	barcodeType = And(Object(map[string]TypeCheck{
		"altText":         IsString,
		"format":          StringEnum(encodings...),
		"message":         IsString,
		"messageEncoding": IsString,
	}), Required("format", "message", "messageEncoding"),
		WhiteList("altText", "format", "message", "messageEncoding"))
)

func init() {
	requiredBaseKeys := Required("description", "formatVersion", "organizationName", "passTypeIdentifier", "serialNumber", "teamIdentifier")
	requiredBaseKeyTypes := Object(map[string]TypeCheck{
		"description":        IsString,
		"formatVersion":      IsNumber,
		"organizationName":   IsString,
		"passTypeIdentifier": IsString,
		"serialNumber":       IsString,
		"teamIdentifier":     IsString,
	})

	assocAppKeyTypes := Object(map[string]TypeCheck{
		"appLaunchURL":               IsString,
		"associatedStoreIdentifiers": ArrayOf(IsNumber),
	})

	companionAppKeyTypes := Object(map[string]TypeCheck{
		"userInfo": Object(map[string]TypeCheck{}),
	})

	expirationKeyTypes := Object(map[string]TypeCheck{
		"expirationDate": isDate,
		"voided":         IsBool,
	})

	relevanceKeyTypes := Object(map[string]TypeCheck{
		"beacons":      ArrayOf(beaconType),
		"locations":    ArrayOf(locationType),
		"maxDistance":  IsNumber,
		"relevantDate": isDate,
	})

	styleType := And(Mutex("boardingPass", "coupon", "eventTicket", "generic", "storeCard"), Object(map[string]TypeCheck{
		"boardingPass": passStructureType,
		"coupon":       passStructureType,
		"eventTicket":  passStructureType,
		"generic":      passStructureType,
		"storeCard":    passStructureType,
	}))

	visualKeyTypes := Object(map[string]TypeCheck{
		"barcode":            barcodeType,
		"backgroundColor":    isColor,
		"foregroundColor":    isColor,
		"groupingIdentifier": IsString,
		"labelColor":         isColor,
		"logoText":           IsString,
		"suppressStripShine": IsBool,
	})

	webServiceKeys := Object(map[string]TypeCheck{
		"authenticationToken": IsString,
		"webServiceURL":       IsString,
	})

	baseList := WhiteList("description", "formatVersion", "organizationName", "passTypeIdentifier", "serialNumber", "teamIdentifier",
		"appLaunchURL", "associatedStoreIdentifiers",
		"userInfo",
		"expriationDate", "voided",
		"beacons", "locations", "maxDistance", "relevantDate",
		"boardingPass", "coupon", "eventTicket", "generic", "storeCard",
		"barcode", "backgroundColor", "foregroundColor", "groupingIdentifier", "labelColor", "logoText", "suppressStripShine",
		"authenticationToken", "webServiceURL")
	checkPass = And(requiredBaseKeys, requiredBaseKeyTypes, assocAppKeyTypes, companionAppKeyTypes, expirationKeyTypes, relevanceKeyTypes, styleType, visualKeyTypes, webServiceKeys, baseList)
}
