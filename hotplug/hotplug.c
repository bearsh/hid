#include <stdlib.h>
#include "hidapi.h"

extern int go_hid_hotplug_cb(hid_hotplug_callback_handle, struct hid_device_info *, hid_hotplug_event, void *);

int hid_hotplug_cb(hid_hotplug_callback_handle callback_handle, struct hid_device_info *device, hid_hotplug_event event, void *user_data) {
	return go_hid_hotplug_cb(callback_handle, device, event, user_data);
}
