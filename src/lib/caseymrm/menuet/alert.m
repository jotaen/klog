#import <Cocoa/Cocoa.h>

#import "alert.h"

void alertClicked(int, const char *);

void showAlert(const char *jsonString) {
	NSDictionary *jsonDict = [NSJSONSerialization
	                          JSONObjectWithData:[[NSString stringWithUTF8String:jsonString]
	                                              dataUsingEncoding:NSUTF8StringEncoding]
	                          options:0
	                          error:nil];

	dispatch_async(dispatch_get_main_queue(), ^{
		NSAlert *alert = [NSAlert new];
		alert.messageText = jsonDict[@"MessageText"];
		alert.informativeText = jsonDict[@"InformativeText"];
		NSArray *buttons = jsonDict[@"Buttons"];
		if (![buttons isEqualTo:NSNull.null] && buttons.count > 0) {
		        for (NSString *label in buttons) {
		                [alert addButtonWithTitle:label];
			}
		}
		NSView *accessoryView;
		NSArray *inputs = jsonDict[@"Inputs"];
		BOOL hasInputs = ![inputs isEqualTo:NSNull.null] && inputs.count > 0;
		if (hasInputs) {
		        BOOL first = false;
		        int y = 30 * inputs.count;
		        accessoryView = [[NSView alloc] initWithFrame:NSMakeRect(0, 0, 200, y)];
		        for (NSString *input in inputs) {
		                y -= 30;
		                NSTextField *textfield =
					[[NSTextField alloc] initWithFrame:NSMakeRect(0, y, 200, 25)];
		                [textfield setPlaceholderString:input];
		                [accessoryView addSubview:textfield];
		                if (!first) {
		                        [alert.window setInitialFirstResponder:textfield];
		                        first = true;
				}
			}
		        [alert setAccessoryView:accessoryView];
		}

		[NSApp activateIgnoringOtherApps:YES];
		NSInteger resp = [alert runModal];
		NSMutableArray *values = [NSMutableArray new];
		if (hasInputs) {
		        for (NSView *subview in accessoryView.subviews) {
		                if (![subview isKindOfClass:[NSTextField class]]) {
		                        continue;
				}
		                [values addObject:((NSTextField *)subview).stringValue];
			}
		}
		NSData *jsonData =
			[NSJSONSerialization dataWithJSONObject:values options:0 error:nil];
		NSString *jsonString =
			[[NSString alloc] initWithData:jsonData encoding:NSUTF8StringEncoding];
		alertClicked(resp - NSAlertFirstButtonReturn, jsonString.UTF8String);
	});
}
