import assert from "node:assert/strict";
import test from "node:test";
import {
  CLASSROOM_TOOL_TABS,
  DEFAULT_CLASSROOM_TOOLS_OPEN,
  DEFAULT_CLASSROOM_TOOL_TAB,
} from "./classroomTools.js";

test("classroom tools expose idle rate and seat query tabs in display order", () => {
  assert.equal(DEFAULT_CLASSROOM_TOOL_TAB, "idle-rate");
  assert.equal(DEFAULT_CLASSROOM_TOOLS_OPEN, false);
  assert.deepEqual(
    CLASSROOM_TOOL_TABS.map((tab) => [tab.id, tab.label]),
    [
      ["idle-rate", "闲置率"],
      ["seat-query", "座位查询"],
    ]
  );
});
