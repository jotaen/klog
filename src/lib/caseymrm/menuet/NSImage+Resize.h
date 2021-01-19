#import <Cocoa/Cocoa.h>

@interface NSImage (Resize)

- (NSImage *)imageWithHeight:(CGFloat)height;

+ (NSImage *)imageFromName:(NSString *)name withHeight:(CGFloat)height;

@end