#ifndef __USERDEFAULTS_H__
#define __USERDEFAULTS_H__
#endif

void setString(const char* key, const char* value);
const char* getString(const char *key);

void setInteger(const char* key, long value);
long getInteger(const char *key);

void setBoolean(const char* key, bool value);
bool getBoolean(const char *key);