import { dayKey, formatDateTime, formatTime } from './datetime';

// These tests run with TZ=UTC (see the "test" script) for stable day boundaries.

describe('dayKey', () => {
  it('groups timestamps that fall on the same local day', () => {
    expect(dayKey('2026-06-24T09:00:00Z')).toBe(dayKey('2026-06-24T14:30:00Z'));
  });

  it('separates timestamps on different days', () => {
    expect(dayKey('2026-06-24T09:00:00Z')).not.toBe(dayKey('2026-06-25T09:00:00Z'));
  });
});

describe('formatters', () => {
  it('format a valid ISO string into a non-empty, time-bearing label', () => {
    expect(formatDateTime('2026-06-24T09:00:00Z')).not.toBe('');
    expect(formatTime('2026-06-24T09:00:00Z')).toMatch(/\d/);
  });

  it('echo invalid input back unchanged', () => {
    expect(formatDateTime('not-a-date')).toBe('not-a-date');
    expect(formatTime('nope')).toBe('nope');
  });
});
