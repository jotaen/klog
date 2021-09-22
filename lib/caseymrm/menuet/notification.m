#import <Cocoa/Cocoa.h>

#import "notification.h"

void showNotification(const char *jsonString) {
	NSDictionary *jsonDict = [NSJSONSerialization
	                          JSONObjectWithData:[[NSString stringWithUTF8String:jsonString]
	                                              dataUsingEncoding:NSUTF8StringEncoding]
	                          options:0
	                          error:nil];
	NSUserNotification *notification = [NSUserNotification new];
	BOOL showsButtons = NO;
	notification.title = jsonDict[@"Title"];
	notification.subtitle = jsonDict[@"Subtitle"];
	notification.informativeText = jsonDict[@"Message"];
	NSString *identifier = jsonDict[@"Identifier"];
	if (identifier.length > 0) {
		notification.identifier = identifier;
	}
	NSString *closeButton = jsonDict[@"CloseButton"];
	if (closeButton.length > 0) {
		showsButtons = true;
		notification.otherButtonTitle = closeButton;
	}
	NSString *actionButton = jsonDict[@"ActionButton"];
	if (actionButton.length > 0) {
		showsButtons = true;
		notification.actionButtonTitle = actionButton;
	}
	NSString *responsePlaceholder = jsonDict[@"ResponsePlaceholder"];
	if (responsePlaceholder.length > 0) {
		notification.hasReplyButton = YES;
		notification.responsePlaceholder = responsePlaceholder;
	}
	if (showsButtons) {
		// Override banner setting, could check plist to see if we're already set to alerts
		[notification setValue:@YES forKey:@"_showsButtons"];
	}
	BOOL removeFromNotificationCenter = [jsonDict[@"RemoveFromNotificationCenter"] boolValue];
	dispatch_async(dispatch_get_main_queue(), ^{
		NSUserNotificationCenter *center =
			[NSUserNotificationCenter defaultUserNotificationCenter];
		[center deliverNotification:notification];
		if (removeFromNotificationCenter) {
		        [center removeDeliveredNotification:notification];
		}
	});
}
