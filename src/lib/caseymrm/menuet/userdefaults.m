#import <Cocoa/Cocoa.h>

#import "menuet.h"

void setString(const char *key, const char *value) {
	NSString *keyStr = [NSString stringWithUTF8String:key];
	NSString *valueStr = [NSString stringWithUTF8String:value];
	[NSUserDefaults.standardUserDefaults setObject:valueStr forKey:keyStr];
	[NSUserDefaults.standardUserDefaults synchronize];
}

const char *getString(const char *key) {
	NSString *keyStr = [NSString stringWithUTF8String:key];
	NSString *valueStr = [NSUserDefaults.standardUserDefaults stringForKey:keyStr];
	return valueStr.UTF8String;
}

void setInteger(const char *key, long value) {
	NSString *keyStr = [NSString stringWithUTF8String:key];
	[NSUserDefaults.standardUserDefaults setInteger:value forKey:keyStr];
	[NSUserDefaults.standardUserDefaults synchronize];
}

long getInteger(const char *key) {
	NSString *keyStr = [NSString stringWithUTF8String:key];
	NSInteger value = [NSUserDefaults.standardUserDefaults integerForKey:keyStr];
	return value;
}

void setBoolean(const char* key, bool value) {
	NSString *keyStr = [NSString stringWithUTF8String:key];
	[NSUserDefaults.standardUserDefaults setBool:value forKey:keyStr];
	[NSUserDefaults.standardUserDefaults synchronize];
}

bool getBoolean(const char *key) {
	NSString *keyStr = [NSString stringWithUTF8String:key];
	NSInteger value = [NSUserDefaults.standardUserDefaults boolForKey:keyStr];
	return value;
}
